package scheduler

import (
	"Gure/gerror"
	"Gure/kits"
	"Gure/logger"
	"Gure/module"
	"errors"
	"fmt"
	"strings"
)

//提供函数支持

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

func (g *gureScheduler) checkPoolsForStart() error {
	//检查各个连接池状态
	if g.reqBuffPool == nil || g.respBuffPool == nil || g.itemBuffPool == nil || g.errBuffPool == nil {
		return fmt.Errorf("failed to initialize the buffer pool")
	}
	if g.reqBuffPool.MaxBufferNum() < 1 {
		return fmt.Errorf("failed to initialize the reqBuffPool")
	}

	if g.respBuffPool.MaxBufferNum() < 1 {
		return fmt.Errorf("failed to initialize the respBuffPool")
	}

	if g.itemBuffPool.MaxBufferNum() < 1 {
		return fmt.Errorf("failed to initialize the itemBuffPool")
	}

	if g.errBuffPool.MaxBufferNum() < 1 {
		return fmt.Errorf("failed to initialize the errBuffPool")
	}
	return nil
}

func (g *gureScheduler) pick() {
	//不断监听数据队列，然后交给数据处理函数去处理
	go func() {
		for true {
			if g.canceled() {
				break //跳循环结束
			}
			data, err := g.itemBuffPool.Get()
			if err != nil {
				logger.Warn("item pool is closed")
				break
			}
			item, ok := data.(module.Item)
			if !ok {
				//数据格式有问题，向error管道发送数据
				errMsg := fmt.Sprintf("incorrect data type %T", item)
				g.sendError(errors.New(errMsg), "")
			}
			//开始执行下载操作
			g.pickOne(item)
		}
	}()
}

func (g *gureScheduler) analyze() {
	go func() {
		for true {
			if g.canceled() {
				break //跳循环结束
			}
			data, err := g.respBuffPool.Get()
			if err != nil {
				logger.Warn("response pool is closed")
				break
			}
			resp, ok := data.(*module.Response)
			if !ok {
				//数据格式有问题，向error管道发送数据
				errMsg := fmt.Sprintf("incorrect data type %T", resp)
				g.sendError(errors.New(errMsg), "")
			}
			//开始执行下载操作
			g.analyzeOne(resp)
		}
	}()

}
func (g *gureScheduler) pickOne(item module.Item) {
	if item == nil {
		return
	}
	//被关闭了
	if g.canceled() {
		return
	}
	ana, err := g.registrar.Get(module.PIPELINE)
	if err != nil {
		g.sendError(fmt.Errorf("couldn't get a pipeline with %s", err), "")
		g.sendData(item)
		return
	}
	pipeline, ok := ana.(module.Pipeline)
	if !ok {
		errMsg := fmt.Sprintf("incorrect data type %T", pipeline)
		g.sendError(errors.New(errMsg), "")
		g.sendData(item)
		return
	}
	errList := pipeline.Send(item)
	if errList != nil {
		for _, err := range errList {
			g.sendError(err, ana.ID())
		}
	}
}

func (g *gureScheduler) analyzeOne(resp *module.Response) {
	if resp == nil {
		return
	}
	//被关闭了
	if g.canceled() {
		return
	}
	ana, err := g.registrar.Get(module.ANALYZER)
	if err != nil {
		g.sendError(fmt.Errorf("couldn't get a downloader with %s", err), "")
		g.sendResp(resp)
		return
	}
	analyzer, ok := ana.(module.Analyzer)
	if !ok {
		errMsg := fmt.Sprintf("incorrect data type %T", analyzer)
		g.sendError(errors.New(errMsg), "")
		g.sendResp(resp)
		return
	}
	dataList, errList := analyzer.Analyze(resp)
	if dataList != nil {
		for _, data := range dataList {
			//检查解析得到的数据
			if data == nil {
				continue
			}
			//数据可能是新的请求
			switch d := data.(type) {
			case *module.Request:
				g.sendReq(d)
			case module.Item:
				g.sendData(d)
			default:
				errMsg := fmt.Sprintf("incorrect data type %T", d)
				g.sendError(errors.New(errMsg), ana.ID())
			}
		}
	}
	if errList != nil {
		for _, err := range errList {
			g.sendError(err, ana.ID())
		}
	}

}

func (g *gureScheduler) download() {
	//开一个goroutine，不断循环读取，一旦cancel就退出就行
	go func() {
		for true {
			if g.canceled() {
				break //跳循环结束
			}
			data, err := g.reqBuffPool.Get()
			if err != nil {
				logger.Warn("request pool is closed")
				break
			}
			request, ok := data.(*module.Request)
			if !ok {
				//数据格式有问题，向error管道发送数据
				errMsg := fmt.Sprintf("incorrect data type %T", request)
				g.sendError(errors.New(errMsg), "")
			}
			//开始执行下载操作
			g.downloadOne(request)
		}
	}()
}

func (g *gureScheduler) downloadOne(request *module.Request) {
	if request == nil { //downloader可能是客户提供的，因此要给出判断
		return
	}
	if g.canceled() {
		return
	}
	get, err := g.registrar.Get(module.DOWNLOADER)
	if err != nil {
		g.sendError(fmt.Errorf("couldn't get a downloader with %s", err), "")
		g.sendReq(request)
		return
	}
	loader, ok := get.(module.DownLoader)
	if !ok {
		g.sendError(fmt.Errorf("incorrect downloader type  %T", loader), "")
		g.sendReq(request)
		return
	}
	resp, err := loader.Download(request)
	//这里才是真正访问过了
	if resp != nil {
		g.sendResp(resp)
	}
	if err != nil {
		g.sendError(err, loader.ID())
	}
	if resp != nil && err == nil {
		g.urlMap.Store(request.HTTPRep().URL.String(), struct{}{})
	}
}

func (g *gureScheduler) canceled() bool {
	select {
	case <-g.ctx.Done():
		return true
	default:
		return false
	}
}

func (g *gureScheduler) sendReq(request *module.Request) bool {
	if request == nil || g.reqBuffPool == nil || g.reqBuffPool.Closed() {
		return false
	}
	//中间需要进行检查操作
	req := request.HTTPRep()
	if req == nil || req.URL == nil {
		return false
	}
	lower := strings.ToLower(req.URL.Scheme)
	if lower != "http" && lower != "https" {
		return false
	}
	//需要检测是否有被访问过
	var ok bool
	_, ok = g.urlMap.Load(req.URL.String())
	if ok {
		return false
	}
	//检查domain,如果在原始链接的domain内
	_, ok = g.acceptedDomain.Load(req.URL.Host)
	if !ok {
		return false
	}
	if request.Depth() > g.maxDepth {
		return false
	}
	go func(resp *module.Request) {
		err := g.reqBuffPool.Put(request)
		if err != nil {
			logger.Warn("request buffer  pool is closed")
		}
	}(request)

	return true

}

func (g *gureScheduler) sendResp(response *module.Response) bool {
	if response == nil || g.respBuffPool == nil || g.respBuffPool.Closed() {
		return false
	}
	go func(resp *module.Response) {
		err := g.respBuffPool.Put(response)
		if err != nil {
			logger.Warn("response buffer  pool is closed")
		}
	}(response)
	return true
}

func (g *gureScheduler) sendData(data module.Data) bool {
	if data == nil || g.itemBuffPool == nil || g.itemBuffPool.Closed() {
		return false
	}
	go func(d module.Data) {
		err := g.itemBuffPool.Put(d)
		if err != nil {
			logger.Warn("item buffer buffer  pool is closed")
		}
	}(data)
	return true
}

func (g *gureScheduler) sendError(err error, mid module.MID) bool {
	//错误发送方法
	if err == nil || g.errBuffPool == nil || g.errBuffPool.Closed() {
		return false //直接发送失败
	}
	//接下来根据mid判断错误类型
	//首先判断是不是框架错误
	spiderError, ok := err.(gerror.SpiderError)
	if !ok { //说明要进行解析操作
		var moduleType module.Type
		var errType module.ErrorType
		//根据mid解析
		spiltMid, err := module.SpiltMid(mid)
		if err != nil {
			//解析失败，说明mid为空调度器error
			errType = module.SchedulerError
		} else {
			//三种类型中的一种
			moduleType = module.LegalLetterTypeMap[spiltMid[0]]
			switch moduleType {
			case module.DOWNLOADER:
				errType = module.DownloaderError
			case module.ANALYZER:
				errType = module.AnalyzerError
			case module.PIPELINE:
				errType = module.PipelineError
			}
		}
		spiderError = gerror.NewSpiderError(errType, err.Error())
	}
	go func(spiderError gerror.SpiderError) {
		if err := g.errBuffPool.Put(spiderError); err != nil {
			logger.Warn("the error pool is closed when put the error")
		}
	}(spiderError)
	return true
}
