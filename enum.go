package vortex

const (
	DissolvedOxygenAndTemperature = iota // 溶氧量
	D2
)

/**
 * 传感器可能出现的状态
 */
const (
	STATUS_NORMAL = iota // 正常运行的
	STATUS_DETACH        // 异常断开的
	STATUS_CLOSED        // 人为关闭的
)
