package scheduler

import (
	"Gure/gerror"
	"Gure/module"
	"fmt"
	"reflect"
)

//Args 提供自检方法
type Args interface {
	Check() error
}

// RequestArgs 请求参数设置
type RequestArgs struct {

	//AcceptedDomains 接受的请求域
	AcceptedDomains []string `json:"acceptedDomains,omitempty"`

	//MaxDepth 最大的请求深度，不允许超过该深度
	MaxDepth uint32 `json:"maxDepth,omitempty"`
}

func (r *RequestArgs) Check() error {
	if r.AcceptedDomains == nil {
		return gerror.NewIllegalParameterError("nil acceptedDomain in reqArgs")
	}
	if r.MaxDepth <= 1 {
		return gerror.NewIllegalParameterError("invalid MaxDepth in reqArgs")
	}
	return nil
}

// DataArgs 数据相关设置
type DataArgs struct {
	ReqBufferCap uint32 `json:"reqBufferCap,omitempty"`

	ReqBufferMaxNum uint32 `json:"reqBufferMaxCap,omitempty"`

	RespBufferCap uint32 `json:"respBufferCap,omitempty"`

	RespBufferMaxNum uint32 `json:"respBufferMaxCap,omitempty"`

	ItemBufferCap uint32 `json:"itemBufferCap,omitempty"`

	ItemBufferMaxNum uint32 `json:"itemBufferMaxCap,omitempty"`

	ErrorBufferCap uint32 `json:"errorBufferCap,omitempty"`

	ErrorBufferMaxNum uint32 `json:"errorBufferMaxCap,omitempty"`
}

// Check 利用反射进行校验
func (r *DataArgs) Check() error {
	va := reflect.ValueOf(*r) //参数必须要是结构体
	for i := 0; i < va.NumField(); i++ {
		field := va.Field(i)
		count := field.Interface()
		if count.(uint32) < 1 {
			return fmt.Errorf("invalid buffer params in dataArgs")
		}
	}
	return nil
}

// ModuleArgs 模块相关设置
type ModuleArgs struct {
	DownLoaders []module.DownLoader
	Analyzers   []module.Analyzer
	Pipelines   []module.Pipeline
}

func (r *ModuleArgs) Check() error {
	//保证每一个都非空，并且长度大于0
	if r.DownLoaders == nil || len(r.DownLoaders) == 0 {
		return fmt.Errorf("invalid module params in moduleArgs")
	}
	if r.Analyzers == nil || len(r.Analyzers) == 0 {
		return fmt.Errorf("invalid module params in moduleArgs")
	}
	if r.Pipelines == nil || len(r.Pipelines) == 0 {
		return fmt.Errorf("invalid module params in moduleArgs")
	}
	return nil
}
