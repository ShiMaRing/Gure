package module

// IDXGenerator 序号生成器接口
type IDXGenerator interface {
	// Start 最小序列号
	Start() uint64
	//Max 最大序列号
	Max() uint64
	//Next 获取下一序列号
	Next() uint64
	//CycleCount 获取循环计数，循环的次数
	CycleCount() uint64
	//Get 获取序列号并准备下一序列号
	Get() uint64
}
