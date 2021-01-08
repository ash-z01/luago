package number

import "strconv"

// 自动类型转换
// 字符串解析为数字

// ParseInteger 字符串解析为整数
func ParseInteger(str string) (int64, bool) {
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err == nil
}

// ParseFloat 字符串解析为浮点数
func ParseFloat(str string) (float64, bool) {
	n, err := strconv.ParseFloat(str, 64)
	return n, err == nil
}
