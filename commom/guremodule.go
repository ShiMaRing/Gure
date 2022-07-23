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
	//TODO implement me
	panic("implement me")
}

func (g gureModule) Addr() string {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) Score() uint64 {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) SetScore(uint642 uint64) {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) ScoreCalculator() module.CalculateScore {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) CalledCount() uint64 {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) AcceptedCount() uint64 {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) CompletedCount() uint64 {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) HandlingNumber() uint64 {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) Counts() module.Counts {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) Summary() module.SummaryStruct {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) IncrCalledCount() {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) IncrAcceptedCount() {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) IncrCompletedCount() {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) IncrHandlingNumber() {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) DecrHandlingNumber() {
	//TODO implement me
	panic("implement me")
}

func (g gureModule) Clear() {
	//TODO implement me
	panic("implement me")
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
