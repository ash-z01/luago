package number

import "math"

// IFloorDiv 整除
func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	}
	return a/b - 1
}

// FFloorDiv 整除
func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

// IMod 取模
// a % b == a - ((a // b) * b)
func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b
}

// FMod 取模
// a % b == a - ((a // b) * b)
func FMod(a, b float64) float64 {
	if a > 0 && math.IsInf(b, 1) || a < 0 && math.IsInf(b, -1) {
		return a
	}
	if a > 0 && math.IsInf(b, -1) || a < 0 && math.IsInf(b, 1) {
		return b
	}
	return a - math.Floor(a/b)*b
}

// ShiftLeft 左移
func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	}
	return ShiftRight(a, -n)
}

// ShiftRight 右移
// Golang中，如果位移运算符左侧是有符号类型，则进行 有符号右移，空位补1，我们期望无符号(空位补0)
func ShiftRight(a, n int64) int64 {
	if n >= 0 {
		return int64(uint64(a) >> uint64(n))
	}
	return ShiftLeft(a, -n)
}

// 自动类型转换

// FloatToInteger 浮点数转为整数
func FloatToInteger(n float64) (int64, bool) {
	i := int64(n)
	return i, float64(i) == n
}
