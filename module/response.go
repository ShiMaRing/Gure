package module

import "net/http"

// Response 数据响应
type Response struct {
	httpResp *http.Response
	depth    uint32
}

func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{httpResp: httpResp, depth: depth}
}

// HTTPRep 获取resp对象
func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

// Depth 获取响应深度
func (resp *Response) Depth() uint32 {
	return resp.depth
}
