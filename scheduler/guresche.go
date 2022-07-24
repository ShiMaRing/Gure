package scheduler

import (
	"Gure/gerror"
	"Gure/kits"
	"Gure/module"
	"Gure/regist"
	"context"
	"fmt"
	"log"
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

func (g *gureScheduler) Start(firstReq *module.Request) error {
	//TODO implement me
	panic("implement me")
}

func (g *gureScheduler) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (g *gureScheduler) Status() Status {
	//TODO implement me
	panic("implement me")
}

func (g *gureScheduler) ErrorChan() <-chan error {
	//TODO implement me
	panic("implement me")
}

func (g *gureScheduler) Idle() bool {
	//TODO implement me
	panic("implement me")
}

func (g *gureScheduler) Summary() SchedulerSummary {
	//TODO implement me
	panic("implement me")
}

//会出现多个线程同时操作调度器并发爬取
func (g *gureScheduler) checkAndSetStatus(wanted Status) (oldStatus Status, err error) {
	//检查当前状态以及期望状态，查看是否能够进行转化操作
	g.statusLock.Lock() //先上锁避免竞态
	defer g.statusLock.Unlock()
	oldStatus = g.status
	err = checkStatus(oldStatus, wanted) //检查是否能够切换
	if err == nil {
		g.status = wanted
		return oldStatus, nil
	}
	return
}

func (g *gureScheduler) setReqArgs(args RequestArgs) {
	//传入的参数设置,默认在检查过程中完成了相应的默认值设置
	for _, domain := range args.AcceptedDomains {
		g.acceptedDomain.Store(domain, struct{}{})
	}
	g.maxDepth = args.MaxDepth
}

func (g *gureScheduler) setDataArgs(args DataArgs) {
	//默认此时参数都已经完成检查了，会设置阈值，少于阈值会进行修订
	g.reqBuffPool = kits.NewPool(args.ReqBufferCap, args.ReqBufferMaxNum)
	g.respBuffPool = kits.NewPool(args.RespBufferCap, args.RespBufferMaxNum)
	g.itemBuffPool = kits.NewPool(args.ItemBufferCap, args.ItemBufferMaxNum)
	g.errBuffPool = kits.NewPool(args.ErrorBufferCap, args.ErrorBufferMaxNum)

}

func (g *gureScheduler) register(args ModuleArgs) (err error) {
	for _, downloader := range args.DownLoaders {
		_, err = g.registrar.Register(downloader)
		if err != nil {
			return
		}
	}
	for _, analyzer := range args.Analyzers {
		_, err = g.registrar.Register(analyzer)
		if err != nil {
			return
		}
	}
	for _, pipeline := range args.Pipelines {
		_, err = g.registrar.Register(pipeline)
		if err != nil {
			return
		}
	}
	return nil
}

func checkStatus(old Status, wanted Status) error {
	//三个状态下不可以改变
	if old == SchedStatusInitializing || old == SchedStatusStopping || old == SchedStatusStarting {
		return gerror.StatusChangeError("current status cant be changed")
	}
	if wanted != SchedStatusInitializing && wanted != SchedStatusStopping && wanted != SchedStatusStarting {
		return gerror.StatusChangeError("wanted status error")
	}
	if old == SchedStatusUninitialized && (wanted == SchedStatusStarting || wanted == SchedStatusStopping) {
		return gerror.StatusChangeError("cant change to wanted")
	}
	if old == SchedStatusStarted && (wanted == SchedStatusStarting || wanted == SchedStatusInitializing) {
		return gerror.StatusChangeError("cant change to wanted")
	}
	if old != SchedStatusStarted && wanted == SchedStatusStopping {
		return gerror.StatusChangeError("cant change to wanted")
	}
	return nil
}

// New 创建原生调度器对象，必须要执行Init方法进行初始化
func New() Scheduler {
	return &gureScheduler{}
}

//NewDefault 生成默认参数调度器对象
func NewDefault() {

}
