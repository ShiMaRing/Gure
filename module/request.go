package module

import "net/http"

//请求的数据类型
type Request struct {
	//http 请求
	httpReq *http.Request
	//爬取深度
	depth uint32
}

func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

// NewRequest 构造方法
func NewRequest(httpReq *http.Request, depth uint32) *Request {
	return &Request{httpReq: httpReq, depth: depth}
}

// HTTPRep 获取req对象
func (req *Request) HTTPRep() *http.Request {
	return req.httpReq
}

// Depth 获取请求深度
func (req *Request) Depth() uint32 {
	return req.depth
}
