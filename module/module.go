package module

type MID string

//mid 的模板，分别为组件类型，序列号，网络地址
var midTmpl = "%s%d|%s"

//Module 为所有基础组件都要实现的基础接口，支持框架扩展
type Module interface {
	//ID 获取当前组件Id
	ID() MID
	//Addr 获取当前网络地址
	Addr() string
	// Score 获取当前组件的评分
	Score() uint64
	//SetScore 设置组件评分
	SetScore(uint642 uint64)
	//ScoreCalculator 评分计算
	ScoreCalculator() CalculateScore
	//CalledCount 被调用次数
	CalledCount() uint64
	//AcceptedCount 最多调用次数
	AcceptedCount() uint64
	//CompletedCount 完成的调用次数
	CompletedCount() uint64
	//HandlingNumber 正在处理的调用数量
	HandlingNumber() uint64
	//Counts 获取所有计数
	Counts() Counts
	//Summary 返回简介
	Summary() SummaryStruct
}

// CalculateScore 评分计算器
type CalculateScore func(counts Counts) uint64

// Counts 所有计数信息
type Counts struct {
	Called uint64 `json:"called,omitempty"`
	//AcceptedCount 最多调用次数
	Accepted uint64 `json:"accepted,omitempty"`
	//CompletedCount 完成的调用次数
	Completed uint64 `json:"completed,omitempty"`
	//HandlingNumber 正在处理的调用数量
	Handling uint64 `json:"handling,omitempty"`
}

//SummaryStruct 摘要
type SummaryStruct struct {
	ID MID `json:"id"`
	Counts
	Extra interface{} `json:"extra"` //提供额外输入
}
