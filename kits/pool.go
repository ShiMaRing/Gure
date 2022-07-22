package kits

// Pool 数据缓冲池，底层由数据缓冲器支撑构造
type Pool interface {
	// BufferCap 返回缓冲池中缓冲器的统一容量
	BufferCap() uint32

	// MaxBufferNum 缓冲池中最大的缓冲容量
	MaxBufferNum() uint32

	// BufferNum 缓冲器数量
	BufferNum() uint32

	// Total 数据总量
	Total() uint64

	// Put 向缓冲池中发送数据，如果池子已经关闭则会报错
	//方法阻塞运行
	Put(data interface{}) error

	// Get 尝试获取数据，关闭会报错
	//方法阻塞运行
	Get() (data interface{}, err error)

	// Close 关闭缓冲池
	//如果关闭成功返回true，关闭失败即之前已经关闭过了返回false
	Close() bool

	// Closed 判断缓冲池是否已经关闭
	Closed() bool
}
