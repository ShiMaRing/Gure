package module

type (
	//Pipeline 条目处理管道接口
	//需要实现并发安全
	Pipeline interface {
		// Module 基础模组
		Module
		// ItemProcessors 所有的处理函数
		ItemProcessors() []ProcessItem
		// Send 发送条目数据进行处理
		Send(item Item) []error
		// FailFast 是否快速失败，一旦某一函数失败则后续全部结束
		FailFast() bool
		// SetFailFast 设置是否快速失败
		SetFailFast(bool)
	}
)

//ProcessItem 处理函数，传入数据，链式传递
type ProcessItem func(item Item) (result Item, err error)
