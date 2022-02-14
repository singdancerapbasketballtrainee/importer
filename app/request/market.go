package request

// 属于是历史遗留问题了，行情目前用的一位市场是一位无符号整数，ifind存的是一位符合整数并且转化成了字符串
// 例如 港股177 存在ifind表里是-79（字符串），但是实际上是177（整数）
func marketUint8toInt8(uMarket uint8) int8 {
	return int8(uMarket)
}
