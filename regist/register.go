package regist

import (
	"Gure/gerror"
	"Gure/module"
	"fmt"
	"sync"
)

//全局注册器
type gureRegister struct {
	//注册map
	moduleTypeMap map[module.Type]map[module.MID]module.Module
	//读写锁保护
	rwLock sync.RWMutex
}

// Register 实例注册方法
func (g *gureRegister) Register(m module.Module) (bool, error) {
	//将module注册
	if m == nil {
		return false, gerror.NewIllegalParameterError("nil module")
	}
	mid := m.ID()
	//parts里面存的是 D P A
	parts, err := module.SpiltMid(mid)
	if err != nil {
		return false, err
	}
	//拿到Type
	moduleType := module.LegalLetterTypeMap[parts[0]]
	//根据拿到的类型和module的具体实现类型进行解析
	if !module.CheckType(moduleType, m) {
		return false, fmt.Errorf("type of module is not equal to %s", moduleType)
	}
	g.rwLock.Lock()
	g.moduleTypeMap[moduleType][mid] = m
	g.rwLock.Unlock()
	return true, nil
}

// Unregister 删除某一注册实例
func (g *gureRegister) Unregister(mid module.MID) (bool, error) {
	//输入mid进行删除操作，先检查mid的合法性
	parts, err := module.SpiltMid(mid)
	if err != nil {
		return false, err
	}
	moduleType := module.LegalLetterTypeMap[parts[0]]
	//使用读写锁避免并发
	g.rwLock.Lock()
	defer g.rwLock.Unlock()
	//查找是否存在该id
	if _, ok := g.moduleTypeMap[moduleType]; !ok {
		return false, fmt.Errorf("no such mid")
	}
	delete(g.moduleTypeMap[moduleType], mid)
	return true, nil
}

// Get 返回得分最低的那个module
func (g *gureRegister) Get(moduleType module.Type) (module.Module, error) {
	//不需要上读写锁，不涉及修改操作
	modules, err := g.GetAllByType(moduleType)
	if err != nil {
		return nil, err
	}
	minScore := uint64(0)            //最小分数
	var selectedModule module.Module //最终选择的结果
	for _, m := range modules {
		err = SetScore(m) //用户自定义的计算方法可能会抛出错误
		if err != nil {
			return nil, err
		}
		score := m.Score()
		if minScore == 0 || score < minScore {
			selectedModule = m
			minScore = score
		}
	}
	return selectedModule, nil
}

// GetAllByType 获取所有注册的该类型的模块
func (g *gureRegister) GetAllByType(moduleType module.Type) (map[module.MID]module.Module, error) {
	g.rwLock.RLock()
	defer g.rwLock.RUnlock()
	//需要检测避免客户包装后传入错误参数
	m, ok := g.moduleTypeMap[moduleType]
	if !ok {
		return nil, fmt.Errorf("unknown moduleType %s ", moduleType)
	}
	return m, nil
}

// GetAll  获取所有的module模块
func (g *gureRegister) GetAll() map[module.MID]module.Module {
	g.rwLock.RLock()
	defer g.rwLock.RUnlock()
	var res = make(map[module.MID]module.Module)
	for _, midMap := range g.moduleTypeMap {
		for mid, value := range midMap {
			res[mid] = value
		}
	}
	return res
}

// Clear 清空所有map
func (g *gureRegister) Clear() {
	g.rwLock.Lock()
	defer g.rwLock.Unlock() //获取锁，避免对错误的map操作

}

// SetScore 为module设置分数，module底层并发安全
func SetScore(m module.Module) error {
	calculator := m.ScoreCalculator()
	if calculator == nil {
		m.SetScore(avgCalculate(m.Counts()))
	} else {
		res, err := calculator(m.Counts())
		if err != nil {
			return err
		}
		//如果有自带方法
		m.SetScore(res)
	}
	return nil
}

func avgCalculate(counts module.Counts) uint64 {
	return (counts.Accepted + counts.Called + counts.Completed + counts.Handling) / 4
}
