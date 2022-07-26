package scheduler

import (
	"Gure/gerror"
	"Gure/kits"
	"Gure/logger"
	"Gure/module"
	"Gure/regist"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type gureMap struct {
	sync.Map
} //并发安全字典

//scheduler 结构体

type gureScheduler struct {
	//最大访问深度
	maxDepth uint32
	//可接受的域名范围
	acceptedDomain gureMap
	//组件注册器
	registrar module.Registrar

	reqBuffPool kits.Pool

	respBuffPool kits.Pool

	itemBuffPool kits.Pool

	errBuffPool kits.Pool

	//完成的链接
	urlMap gureMap

	//用于停止调度器
	ctx context.Context

	//停止调度器
	cancelFunc context.CancelFunc

	//状态
	status Status

	//状态锁
	statusLock sync.RWMutex

	//摘要
	summary SchedulerSummary
}

func (g *gureScheduler) Init(reqArgs RequestArgs, dataArgs DataArgs, moduleArgs ModuleArgs) (err error) {
	//首先需要检查调度器状态
	var oldStatus Status //一开始是未初始化状态
	//检查通过后会设置为想要得到的状态
	oldStatus, err = g.checkAndSetStatus(SchedStatusInitializing)
	if err != nil {
		return err
	}
	//最后的状态切换，保证转换切换正确
	defer func() {
		g.statusLock.Lock()
		if err != nil {
			//rollback
			g.status = oldStatus
		} else {
			g.status = StatusInitialized
		}
		g.statusLock.Unlock()
	}()
	//状态切换完成，开始初始化,检查各个设置状态
	if err = reqArgs.Check(); err != nil {
		log.Println("reqArgs check fail")
		return err
	}
	if err = dataArgs.Check(); err != nil {
		log.Println("dataArgs check fail")
		return err
	}
	if err = moduleArgs.Check(); err != nil {
		log.Println("moduleArgs check fail")
		return err
	}

	//三个都校验完成，开始初始化操作
	//初始化链接保存map
	g.acceptedDomain = gureMap{}
	g.urlMap = gureMap{}

	//初始化取消上下文
	ctx, cancelFunc := context.WithCancel(context.Background())
	g.ctx = ctx
	g.cancelFunc = cancelFunc
	//初始化注册器
	g.registrar = regist.NewRegister()

	err = g.register(moduleArgs)
	err = fmt.Errorf("register module fail with %v", err)

	g.setReqArgs(reqArgs)
	g.setDataArgs(dataArgs)
	//注册
	g.summary = &SummaryStruct{
		RequestArgs: reqArgs,
		DataArgs:    dataArgs,
		ModuleArgs:  moduleArgs,
	}

	//完成之后生成最后的摘要结构体,可能会出现注册失败场景
	return err
}

func (g *gureScheduler) Start(firstReq *http.Request) (err error) {
	//捕获panic防止程序崩溃
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Schedule error %s ", p)
			logger.Fatalf(errMsg)
			err = errors.New(errMsg)
		}
	}()
	logger.Info("Start scheduler...")

	//首先需要检查状态是否合法
	var oldStatus Status
	oldStatus, err = g.checkAndSetStatus(SchedStatusStarting)
	//出现错误应该要换回原来的状态,否则认为启动完成
	defer func() {
		g.statusLock.Lock()
		if err != nil {
			g.status = oldStatus
		} else {
			g.status = SchedStatusStarted
		}
		g.statusLock.Unlock()
	}()
	//查看状态是否切换成功，失败则直接退出
	if err != nil {
		return
	}
	//检查传入的初始参数
	if firstReq == nil {
		err = gerror.NewIllegalParameterError("nil firstReq")
		return
	}
	//获得初次的域名并进行添加
	var primaryDomain string
	if firstReq.Host == "" {
		err = gerror.NewIllegalParameterError("empty host")
		return
	}
	primaryDomain = firstReq.Host
	g.acceptedDomain.Store(primaryDomain, struct{}{})
	//开始执行各个操作，还需要检查缓冲池的初始化问题
	if err = g.checkPoolsForStart(); err != nil {
		return err
	}
	//都没有问题就可以开始爬取工作
	//三个异步方法，同步爬取
	g.download()
	g.analyze()
	g.pick()
	request := module.NewRequest(firstReq, 0)
	g.sendReq(request) //向缓冲池放入第一个请求
	return nil
}

func (g *gureScheduler) Stop() (err error) {
	//捕获panic防止程序崩溃
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Schedule close with panic %s ", p)
			logger.Fatalf(errMsg)
			err = errors.New(errMsg)
		}
	}()
	logger.Info("Closing scheduler...")

	//首先需要检查状态是否合法
	var oldStatus Status
	oldStatus, err = g.checkAndSetStatus(SchedStatusStopping)
	//出现错误应该要换回原来的状态,否则认为启动完成
	defer func() {
		g.statusLock.Lock()
		if err != nil {
			g.status = oldStatus
		} else {
			g.status = SchedStatusStopped
		}
		g.statusLock.Unlock()
	}()
	//查看状态是否切换成功，失败则直接退出
	if err != nil {
		return
	}
	//调用ctx，调用close,应当先取消，避免新请求被处理
	g.cancelFunc()
	g.reqBuffPool.Close()
	g.respBuffPool.Close()
	g.itemBuffPool.Close()
	g.errBuffPool.Close()
	logger.Info("finish close scheduler")
	return
}

func (g *gureScheduler) Status() Status {
	g.statusLock.RLock()
	defer g.statusLock.RUnlock()
	return g.status
}

// ErrorChan 返回只读错误通道，用户不应当主动关闭
func (g *gureScheduler) ErrorChan() (<-chan error, error) {
	errBuffer := g.errBuffPool
	if errBuffer == nil {
		return nil, fmt.Errorf("the errBufferPool is nil")
	}
	errCh := make(chan error, errBuffer.BufferCap())
	go func(errBuffer kits.Pool, errCh chan error) {
		for {
			if g.canceled() {
				close(errCh) //关闭通道，外界不应该主动关闭
				return
			}
			data, err := errBuffer.Get()
			if err != nil {
				//表示连接池关闭
				logger.Warn("errBufferPool is closed")
				close(errCh)
				return
			}
			err, ok := data.(error)
			if !ok {
				errMsg := fmt.Sprintf("incorrect error type %T", err)
				logger.Warn(errMsg)
				g.sendError(errors.New(errMsg), "")
			}
			if g.canceled() {
				close(errCh)
				return
			} else {

			}

		}

	}(errBuffer, errCh)

	return errCh, nil
}

func (g *gureScheduler) Idle() bool {
	//观察是否空闲，遍历所有模块
	var idle bool
	for _, m := range g.registrar.GetAll() {
		if m.HandlingNumber() > 0 {
			idle = false
			return idle
		}
	}
	return true
}

// Summary 返回当前的结构体状态
func (g *gureScheduler) Summary() SchedulerSummary {
	temp := g.summary //存储了配置信息
	summaryStruct := temp.(*SummaryStruct)
	summaryStruct.Status = g.Status()
	summaryStruct.Downloaders = []module.SummaryStruct{}
	summaryStruct.Pipelines = []module.SummaryStruct{}
	summaryStruct.Analyzers = []module.SummaryStruct{}
	byType, _ := g.registrar.GetAllByType(module.DOWNLOADER)
	for _, m := range byType {
		summaryStruct.Downloaders = append(summaryStruct.Downloaders, convertToSummary(m))
	}

	byType, _ = g.registrar.GetAllByType(module.ANALYZER)
	for _, m := range byType {
		summaryStruct.Analyzers = append(summaryStruct.Analyzers, convertToSummary(m))
	}

	byType, _ = g.registrar.GetAllByType(module.PIPELINE)
	for _, m := range byType {
		summaryStruct.Pipelines = append(summaryStruct.Pipelines, convertToSummary(m))
	}

	return summaryStruct
}

func convertToSummary(m module.Module) module.SummaryStruct {
	var summaryStruct module.SummaryStruct
	summaryStruct.ID = m.ID()
	summaryStruct.Completed = m.CompletedCount()
	summaryStruct.Handling = m.HandlingNumber()
	summaryStruct.Called = m.CalledCount()
	return summaryStruct
}
