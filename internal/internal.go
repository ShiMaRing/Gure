package internal

import "Gure/module"

// ModuleInternal 该接口实现组件的通用功能，
type ModuleInternal interface {
	module.Module
	// IncrCalledCount 增加调用计数
	IncrCalledCount()
	// IncrAcceptedCount 增加接受计数
	IncrAcceptedCount()
	// IncrCompletedCount 增加成功计数
	IncrCompletedCount()
	// IncrHandlingNumber 增加实时处理计数
	IncrHandlingNumber()
	// DecrHandlingNumber 把实时计数减一
	DecrHandlingNumber()
	// Clear 清空计数
	Clear()
}
