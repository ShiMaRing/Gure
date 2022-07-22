package scheduler

type Status int

const (
	// SchedStatusUninitialized 未初始化的状态
	SchedStatusUninitialized Status = iota
	// SchedStatusInitializing 正在初始化的状态
	SchedStatusInitializing
	// StatusInitialized 巳初始化的状态
	StatusInitialized
	// SchedStatusStarting 正在启动的状态
	SchedStatusStarting
	// SchedStatusStarted 已启动的状态
	SchedStatusStarted
	// SchedStatusStopping 正在停止的状态
	SchedStatusStopping
	// SchedStatusStopped 已停止的状态
	SchedStatusStopped
)
