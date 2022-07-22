package module

//组件注册器的接口
type Registrar interface {
	// Register 用于注册组件实例
	Register(module Module) (bool, error)

	// Unregister 用于注销组件实例
	Unregister(mid MID) (bool, error)

	// Get 用于获取一个指定类型的组件的实例，
	//该函数基于负载均衡策略返回实例
	Get(moduleType Type) (Module, error)

	// GetAllByType 用于获取指定类型的所有组件实例
	GetAllByType(moduleType Type) (map[MID]Module, error)

	// GetAll 用于获取所有组件实例
	GetAll() map[MID]Module

	// Clear 清除所有的组件注册记录
	Clear()
}
