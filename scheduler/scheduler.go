package scheduler

import "Gure/module"

type Scheduler interface {
	// Init 初始化参数
	Init(args RequestArgs, dataArgs DataArgs, moduleArgs ModuleArgs)

	// Start 开始爬取第一个请求
	Start(firstReq *module.Request) error

	// Stop 暂停当前爬取工作
	Stop() error

	// Status 返回当前状态
	Status() Status

	// ErrorChan 返回通道，接收过程中的错误
	ErrorChan() <-chan error

	// Idle 返回当前模块的空闲状态
	Idle() bool

	// Summary 返回调度器摘要
	Summary() SchedulerSummary
}
