package module

// DownLoader 下载器模块
//要求并发安全
type DownLoader interface {
	//包含基础组件Module
	Module
	// Download 下载方法
	Download(req *Request) (*Response, error)
}
