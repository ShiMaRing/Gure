package kits

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ParameterIllegalError = errors.New("parameterIllegal error")
var BufferClosedError = errors.New("buffer closed error")

//Buffer 缓冲器,缓冲池的底层数据类型，扩展通道数据结构
type Buffer interface {

	// Cap 返回缓冲池中缓冲器的统一容量
	Cap() uint32

	// Len 缓冲池中最大的缓冲容量
	Len() uint32

	// Put 向缓冲池中发送数据，如果池子已经关闭则会报错
	//方法阻塞运行
	Put(data interface{}) (bool, error)

	// Get 尝试获取数据，关闭会报错
	//方法阻塞运行
	Get() (data interface{}, err error)

	// Close 关闭缓冲池
	//如果关闭成功返回true，关闭失败即之前已经关闭过了返回false
	Close() bool

	// Closed 判断缓冲池是否已经关闭
	Closed() bool
}

type gureBuffer struct {
	ch          chan interface{} //数据存放通道
	closed      int32            //缓冲器的关闭状态
	closingLock sync.RWMutex     //读写锁，避免竞态
}

func NewGureBuffer(size uint32) (*gureBuffer, error) {
	if size == 0 {
		return nil, ParameterIllegalError
	}
	return &gureBuffer{
		ch: make(chan interface{}, size),
	}, nil
}

func (g *gureBuffer) Cap() uint32 {
	return uint32(cap(g.ch))
}

func (g *gureBuffer) Len() uint32 {
	return uint32(len(g.ch))
}

// Put 非阻塞方法，如果无法放入就返回false，向关闭的通道发送数据会出错
func (g *gureBuffer) Put(data interface{}) (bool, error) {
	//先尝试获取锁
	g.closingLock.RLock() //尝试获取读锁
	defer g.closingLock.RUnlock()
	//写入获取都是获取读锁，关闭操作必须要获取写锁，因为是互斥的，关闭操作必须要发生在读取写入完成之后
	if g.Close() {
		return false, BufferClosedError
	}
	select {
	case g.ch <- data:
		return true, nil
	default:
		return false, nil
	}
}

// Get 非阻塞方法，如果向关闭的通道接收数据将会返回错误
func (g *gureBuffer) Get() (interface{}, error) {
	select {
	case data, ok := <-g.ch:
		if !ok {
			return nil, BufferClosedError
		}
		return data, nil
	default:
		return nil, nil
	}
}

func (g *gureBuffer) Close() bool {
	//原子操作，如果为 0 就更换为 1，使用原子操作避免重复关闭
	if atomic.CompareAndSwapInt32(&g.closed, 0, 1) {
		//需要获取写锁，与读锁互斥
		g.closingLock.Lock()
		close(g.ch)
		g.closingLock.Unlock()
		return true
	}
	return false
}

// Closed 检查是否已经关闭
func (g *gureBuffer) Closed() bool {
	if atomic.LoadInt32(&g.closed) == 0 {
		return false
	}
	return true
}
