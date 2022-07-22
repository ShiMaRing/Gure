package scheduler

import "Gure/module"

// RequestArgs 请求参数设置
type RequestArgs struct {

	//AcceptedDomains 接受的请求域
	AcceptedDomains []string `json:"acceptedDomains,omitempty"`

	//MaxDepth 最大的请求深度，不允许超过该深度
	MaxDepth uint32 `json:"maxDepth,omitempty"`
}

// DataArgs 数据相关设置
type DataArgs struct {
	ReqBufferCap uint32 `json:"reqBufferCap,omitempty"`

	ReqBufferMaxCap uint32 `json:"reqBufferMaxCap,omitempty"`

	RespBufferCap uint32 `json:"respBufferCap,omitempty"`

	RespBufferMaxCap uint32 `json:"respBufferMaxCap,omitempty"`

	ItemBufferCap uint32 `json:"itemBufferCap,omitempty"`

	ItemBufferMaxCap uint32 `json:"itemBufferMaxCap,omitempty"`

	ErrorBufferCap uint32 `json:"errorBufferCap,omitempty"`

	ErrorBufferMaxCap uint32 `json:"errorBufferMaxCap,omitempty"`
}

// ModuleArgs 模块相关设置
type ModuleArgs struct {
	DownLoaders []module.DownLoader
	Analyzers   []module.Analyzer
	Pipelines   []module.Pipeline
}

type Args interface {
	Check() error
}
