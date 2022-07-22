package scheduler

import "Gure/module"

// SchedulerSummary 调度器摘要接口
type SchedulerSummary interface {
	// Struct 获取摘要的结构体信息
	Struct() SummaryStruct
	// String 获取摘要的字符串信息
	String() string
}

// SummaryStruct 调度器摘要结构体，由调度器获取参数并构造返回
type SummaryStruct struct {
	RequestArgs RequestArgs
	DataArgs    DataArgs
	ModuleArgs  ModuleArgs
	Status      Status
	Downloaders []module.SummaryStruct
	Analyzers   []module.SummaryStruct
	Pipelines   []module.SummaryStruct
	NumUrl      uint64
}
