package commom

import (
	"Gure/internal"
	"Gure/module"
	"fmt"
)

//实现ModuleInternal，内嵌基本组件
type gureModule struct {
	//组件id
	mid module.MID
	//网络地址
	addr string
	//评分
	score uint64
	//评分计算器
	scoreCalculator module.CalculateScore
	//被调用次数
	calledCount uint64
	//最多接收次数
	acceptedCount uint64
	//成功完成次数
	completedCount uint64
	//正在处理的次数
	handlingNumber uint64
}

func (g gureModule) ID() module.MID {
	return g.mid
}

func (g gureModule) Addr() string {
	return g.addr
}

func (g gureModule) Score() uint64 {
	return g.score
}

func (g gureModule) SetScore(uint642 uint64) {
	g.score = uint642
}

func (g gureModule) ScoreCalculator() module.CalculateScore {
	return g.scoreCalculator
}

func (g gureModule) CalledCount() uint64 {
	return g.calledCount
}

func (g gureModule) AcceptedCount() uint64 {
	return g.acceptedCount
}

func (g gureModule) CompletedCount() uint64 {
	return g.completedCount
}

func (g gureModule) HandlingNumber() uint64 {
	return g.handlingNumber
}

func (g gureModule) Counts() module.Counts {
	return g.Counts()
}

func (g gureModule) Summary() module.SummaryStruct {
	return g.Summary()
}

func (g gureModule) IncrCalledCount() {
	g.calledCount++
}

func (g gureModule) IncrAcceptedCount() {
	g.acceptedCount++
}

func (g gureModule) IncrCompletedCount() {
	g.completedCount++
}

func (g gureModule) IncrHandlingNumber() {
	g.handlingNumber++
}

func (g gureModule) DecrHandlingNumber() {
	g.handlingNumber--
}

func (g gureModule) Clear() {
	g.acceptedCount = 0
	g.completedCount = 0
	g.calledCount = 0
	g.handlingNumber = 0
}

func NewModuleInternal(mid module.MID, scoreCalculator module.CalculateScore) (internal.ModuleInternal, error) {
	//分割mid,检查id是否符合规范
	parts, err := module.SpiltMid(mid)
	if err != nil {
		return nil, fmt.Errorf("wrong mid with %s", mid)
	}
	return &gureModule{mid: mid,
		scoreCalculator: scoreCalculator, addr: parts[2]}, nil
}
