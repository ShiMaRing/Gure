package kits

import (
	"sync"
	"sync/atomic"
)

const EXTRA = 10

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
type gurePool struct {
	//统一缓冲器大小
	bufferCap uint32
	//最大缓冲器数量
	maxBufferNum uint32
	//缓冲器实际数量
	bufferNumber uint32
	//池中数据总数
	total uint64
	//存放数据的通道
	bufChan chan *gureBuffer
	//缓冲池状态
	closed uint32
	//读写保护
	rwLock sync.RWMutex
}

func (g *gurePool) BufferCap() uint32 {
	return g.bufferCap
}

func (g *gurePool) MaxBufferNum() uint32 {
	return g.maxBufferNum
}

// BufferNum 原子操作，避免竞态
func (g *gurePool) BufferNum() (num uint32) {
	return atomic.LoadUint32(&g.bufferNumber)
}

func (g *gurePool) Total() uint64 {
	return uint64(atomic.LoadUint64(&g.total))
}

func (g *gurePool) Put(data interface{}) (err error) {
	//取出一个buffer填充后重新放入
	if g.Closed() {
		return BufferClosedError
	}
	var count uint32              //记录失败次数
	maxCount := g.BufferNum() * 3 //当前数量的五倍
	var ok bool
	for buf := range g.bufChan {
		ok, err = g.putData(buf, data, &count, maxCount)
		if ok || err != nil { //数据填入成功或者出现错误，即pool关闭就返回
			break
		}
	}
	return
}

func (g *gurePool) Get() (data interface{}, err error) {
	//判断是否已经关闭
	if g.Closed() {
		return
	}
	var count uint32
	maxCount := g.BufferNum() * 8 //当前数量的五倍
	for buf := range g.bufChan {
		data, err = g.getData(buf, &count, maxCount)
		if data != nil || err != nil {
			break
		}
	}
	return
}

func (g *gurePool) Close() bool {
	if g.Closed() {
		return false
	}
	//尝试获取锁,写锁与读锁互斥
	g.rwLock.Lock()
	defer g.rwLock.Unlock()
	close(g.bufChan) //关闭所有的buf
	for buf := range g.bufChan {
		buf.Close()
	}
	return true
}

func (g *gurePool) Closed() bool {
	return atomic.LoadUint32(&g.closed) == 1
}

func (g *gurePool) putData(buf *gureBuffer, data interface{}, count *uint32, max uint32) (ok bool, err error) {
	//原子判断是否关闭
	if g.Closed() {
		return false, BufferClosedError
	}
	defer func() {
		g.rwLock.RLock() //获取读锁，避免关闭操作发生
		//检查是否关闭
		if g.Closed() { //缓存数量减一，表示没有重新置入。原子操作避免竞态
			atomic.AddUint32(&g.bufferNumber, ^uint32(0))
			err = BufferClosedError
			buf.Close() //及时关闭避免内存泄露
		} else {
			//重新加回去，此时是线程安全的，这个步骤可能会产生阻塞，因为此时通道可能已经满了
			g.bufChan <- buf
		}
		g.rwLock.RUnlock() //解锁
	}()
	//如果没有关闭,返回false表示没有写入失败，err不为nil表示通道关闭
	ok, err = buf.Put(data)
	if ok { //写入成功，直接返回即可,此时添加total，表示数据量+1
		atomic.AddUint64(&g.total, 1)
		return
	}
	//出现错误，表示通道关闭,直接返回即可，get还要用到
	if err != nil {
		return
	}
	//此时是没有err也没有放入成功，需要进行false操作,count操作是线程安全的
	*count++
	//count表示向所有的buffer放置数据
	//判断是否需要新增buffer,需要调用原子方法避免竞态操作
	if *count >= max && g.BufferNum() < g.MaxBufferNum() {
		//说明此时能够创建新的
		//尝试获取锁,不因该是读写锁，避免新增的时候被填满了
		g.rwLock.Lock()
		//避免获取锁的过程中被填满了或者被关闭了
		if g.BufferNum() < g.MaxBufferNum() {
			if g.Closed() {
				g.rwLock.Unlock()
				return
			}
			newBuf, _ := NewGureBuffer(g.bufferCap)
			newBuf.Put(data)
			g.bufChan <- newBuf
			atomic.AddUint32(&g.bufferNumber, 1)
			atomic.AddUint64(&g.total, 1)
			ok = true
		}
		//如果已经被填满了，解锁返回
		g.rwLock.Unlock()
		return

	}
	return
}

func (g *gurePool) getData(buf *gureBuffer, count *uint32, max uint32) (data interface{}, err error) {
	//大致类似
	if g.Closed() {
		return nil, BufferClosedError
	}
	defer func() {
		//如果多次获取仍然获取不到，并且当前缓冲区没有数据，缓冲池还有其他缓冲器
		if *count >= max && buf.Len() == 0 && g.BufferNum() > 1 {
			//缩小缓冲区
			buf.Close()
			atomic.AddUint32(&g.bufferNumber, ^uint32(0))
			*count = 0 //重新计数
			return
		}
		//否则的话就塞回去，塞回去要获得读写锁，避免中途被关闭了
		g.rwLock.RLock()
		if g.Closed() {
			//连接池关闭，不塞回去，关闭通道避免内存泄漏
			atomic.AddUint32(&g.bufferNumber, ^uint32(0))
			err = BufferClosedError //错误赋值
			buf.Close()             //及时关闭避免内存泄露
		} else {
			//塞回去
			g.bufChan <- buf
		}
		g.rwLock.RUnlock()
	}()
	//尝试获取数据，此时buf只有一个线程操作，是安全的,出现错误则是因为通道被关闭
	data, err = buf.Get()
	if data != nil {
		atomic.AddUint64(&g.total, ^uint64(0))
		return data, nil
	}
	if err != nil {
		return
	}
	*count++ //拿不到数据
	return
}

func NewPool(bufferCap, bufferMaxNum uint32) Pool {

	var gure = &gurePool{}
	gure.bufferCap = bufferCap
	gure.maxBufferNum = bufferMaxNum
	//通道缓冲数，额外添加一部分区域，用来减少阻塞
	gure.bufChan = make(chan *gureBuffer, bufferMaxNum+EXTRA)
	buffer, _ := NewGureBuffer(gure.bufferCap)
	gure.bufChan <- buffer
	return gure
}
