package module

type Analyzer interface {
	// Module 基础模型
	Module
	//RespParsers 返回所有的解析方法
	RespParsers() []ParseResponse
	// Analyzer 经过解析方法，得到最终的数据结果
	Analyzer(resp *Response) ([]Data, error)
}

// ParseResponse 解析响应的函数类型
type ParseResponse func(resp *Response, respDepth uint32) ([]Data, error)
