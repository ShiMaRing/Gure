package downloader

import (
	"Gure/commom"
	"Gure/gerror"
	"Gure/internal"
	"Gure/module"
	"log"
	"net/http"
)

//内置下载器,内嵌结构体实现接口
type gureDownloader struct {
	internal.ModuleInternal             //基础方法交给基础实例去实现即可
	httpClient              http.Client //http客户端提供下载方法
}

func (g *gureDownloader) Download(req *module.Request) (*module.Response, error) {
	//处理计数增加，完成后处理计数减少
	g.IncrHandlingNumber()
	defer g.DecrHandlingNumber()
	g.IncrCalledCount()
	if req == nil {
		return nil, gerror.NewIllegalParameterError("nil req")
	}
	httpReq := req.HTTPRep() //返回封装的req
	if httpReq == nil {
		return nil, gerror.NewIllegalParameterError("nil httpReq")
	}
	//允许开始工作
	g.IncrAcceptedCount()
	log.Printf("Do DownLoader : URL %s Depth %d ...\n", httpReq.URL, req.Depth())
	res, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	g.IncrCompletedCount()
	return module.NewResponse(res, req.Depth()), nil
}

// New 应该返回接口类型，为扩展做好准备
//命名技巧，当独占一个包的时候可以省略表示名
func New(mid module.MID, client *http.Client, scoreCalculator module.CalculateScore) (module.DownLoader, error) {
	moduleInternal, err := commom.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, gerror.NewIllegalParameterError("nil http client")
	}
	return &gureDownloader{
		ModuleInternal: moduleInternal,
		httpClient:     *client,
	}, nil
}
