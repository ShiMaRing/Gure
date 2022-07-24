package analyzer

import (
	"Gure/commom"
	"Gure/gerror"
	"Gure/internal"
	"Gure/kits"
	"Gure/module"
	"log"
)

type gureAnalyzer struct {
	internal.ModuleInternal
	respParsers []module.ParseResponse
}

func (g *gureAnalyzer) RespParsers() []module.ParseResponse {
	if g.respParsers == nil {
		g.respParsers = make([]module.ParseResponse, 0)
	}
	return g.respParsers
}

// Analyze 解析方法
func (g *gureAnalyzer) Analyze(resp *module.Response) (dataList []module.Data, errorList []error) {
	//首先需要保证数据的正确性完整性
	g.IncrHandlingNumber()
	defer g.DecrHandlingNumber()
	g.IncrCalledCount()
	if resp == nil {
		errorList = append(errorList, gerror.NewIllegalParameterError("nil resp "))
		return
	}

	httpRes := resp.HTTPResp()

	if httpRes == nil {
		errorList = append(errorList, gerror.NewIllegalParameterError("nil httpResp "))
		return
	}
	request := httpRes.Request
	if httpRes == nil {
		errorList = append(errorList, gerror.NewIllegalParameterError("nil request "))
		return
	}
	url := request.URL
	if url == nil {
		errorList = append(errorList, gerror.NewIllegalParameterError("nil url "))
		return
	}
	g.IncrAcceptedCount() //接受数量提高
	depth := resp.Depth()
	log.Printf("Parse the response Url: %s depth: %d \n", url, depth)

	//下面开始解析工作
	if httpRes.Body != nil {
		defer httpRes.Body.Close() //及时关闭链接
	}
	multipleReader, err := kits.NewMultipleReader(httpRes.Body)
	if err != nil {
		return nil, append(errorList, err)
	}
	//应当保证不为nil，提前准备部分缓冲区提高效率
	dataList = make([]module.Data, len(g.respParsers))
	for _, respParse := range g.RespParsers() {
		httpRes.Body = multipleReader.Reader()                 //将数据流转换为新的readercloser，提供重复读取功能
		parseList, errList := respParse(httpRes, resp.Depth()) //解析得到相关数据
		if parseList != nil {                                  //这里是用户传入的方法，不可以信任
			for _, value := range parseList {
				if value == nil {
					continue
				}
				dataList = append(dataList, value) //添加到结果列表
			}
		}
		if errList != nil { //这里是用户传入的方法，不可以信任
			for _, value := range errList {
				if value == nil {
					continue
				}
				errList = append(errorList, value) //添加到结果列表
			}
		}
	}
	//完成解析，判断是否完全完成，不出现错误
	if len(errorList) == 0 {
		g.IncrCompletedCount()
	}
	return dataList, errorList
}

//返回一个分析器，参数需要解析方法
func New(mid module.MID, respParsers []module.ParseResponse, scoreCalculator module.CalculateScore) (module.Analyzer, error) {
	moduleInternal, err := commom.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if respParsers == nil {
		return nil, gerror.NewIllegalParameterError("nil resParsers")
	}
	if len(respParsers) == 0 {
		return nil, gerror.NewIllegalParameterError("empty resParses")
	}
	var innerParsers []module.ParseResponse
	//检查外来代码
	for _, f := range respParsers {
		if f == nil {
			return nil, gerror.NewIllegalParameterError("nil resParser func")
		}
		innerParsers = append(innerParsers, f)
	}
	return &gureAnalyzer{
		ModuleInternal: moduleInternal,
		respParsers:    innerParsers,
	}, nil

}
