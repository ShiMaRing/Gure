package pipeline

import (
	"Gure/commom"
	"Gure/gerror"
	"Gure/internal"
	"Gure/module"
)

//条目处理管道，传递来的数据通过管道传递处理
type gurePipeline struct {
	//组件基础实例
	internal.ModuleInternal
	//条目处理器，传入参数进行解析
	itemProcessors []module.ProcessItem
	//是否需要快速失败
	failFast bool
}

func (g *gurePipeline) ItemProcessors() []module.ProcessItem {
	return g.itemProcessors
}

// Send 发送操作，需要将一系列数据进行处理
func (g *gurePipeline) Send(item module.Item) []error {
	g.IncrCalledCount()
	g.IncrHandlingNumber()
	defer g.DecrHandlingNumber()

	//检查内容
	var errList []error
	//检查item
	if item == nil || len(item) == 0 {
		errList = append(errList, gerror.NewIllegalParameterError("invalid item"))
		return errList
	}
	g.IncrAcceptedCount()
	//内容无误开始发送
	var temp module.Item = item
	for _, f := range g.ItemProcessors() {
		next, err := f(temp)
		if err != nil {
			if g.FailFast() { //直接退出即可
				errList = append(errList, err)
				return errList
			}
			errList = append(errList, err)
		}
		if next != nil {
			temp = next
		}
	}
	if len(errList) == 0 {
		g.IncrCompletedCount()
	}
	return errList
}

func (g *gurePipeline) FailFast() bool {
	return g.failFast
}

func (g *gurePipeline) SetFailFast(b bool) {
	g.failFast = b
}

func New(mid module.MID, scoreCalculator module.CalculateScore, itemProcessors []module.ProcessItem, fastFail bool) (module.Pipeline, error) {
	moduleInternal, err := commom.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if itemProcessors == nil {
		return nil, gerror.NewIllegalParameterError("nil processors")
	}
	if len(itemProcessors) == 0 {
		return nil, gerror.NewIllegalParameterError("empty processors")
	}

	return &gurePipeline{
		ModuleInternal: moduleInternal,
		itemProcessors: itemProcessors,
		failFast:       fastFail,
	}, nil
}
