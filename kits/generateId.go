package kits

import "sync"

type IdGenerator struct {
	StartNum   uint64
	currentNum uint64
	maxNum     uint64
	cycle      uint64
	lock       sync.RWMutex
}

func NewIdGenerator(start, max int) *IdGenerator {
	return &IdGenerator{
		StartNum:   uint64(start),
		currentNum: uint64(start),
		maxNum:     uint64(max),
		lock:       sync.RWMutex{},
	}

}

func (i *IdGenerator) Start() uint64 {
	//返回最小值
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.StartNum
}

func (i *IdGenerator) Max() uint64 {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.maxNum
}

func (i *IdGenerator) Next() uint64 {
	i.lock.Lock()
	defer i.lock.Unlock()
	if i.currentNum == i.maxNum {
		i.currentNum = i.StartNum
		i.cycle++
	} else {
		i.currentNum++
	}
	return i.currentNum
}

func (i *IdGenerator) CycleCount() uint64 {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.cycle
}

func (i *IdGenerator) Get() uint64 {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.currentNum
}
