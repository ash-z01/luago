package vm

/*
[iABC模式]中,操作数 B 和 C 只有9个比特，如果当作无符号整数的话，最大也不能超过512
但因为表构造器便捷实用，所以Lua也经常被用来描述数据[类似JSON]
如果有很大的数据需要写成表构造器，但是表的初始容量又不够大，就容易导致表频繁扩容从而影响数据加载效率
为了解决这个问题，NEWTABLE指令的B和C操作数使用了一种叫作浮点字节（Floating Point Byte）的编码方式
这种编码方式和浮点数的编码方式类似，只是仅用一个字节
具体来说，如果把某个字节用二进制写成eeeeexxx，
那么当 eeeee == 0 时，该字节表示的整数就是xxx，否则该字节表示的整数是 (1xxx) ＊ 2^(eeeee - 1)
*/

// Int2fb convert an integer to "floating point byte"
func Int2fb(x int) int {
	e := 0 /* exponent */
	if x < 8 {
		return x
	}
	for x >= (8 << 4) { /* coarse steps */
		x = (x + 0xf) >> 4 /* x = ceil(x / 16) */
		e += 4
	}
	for x >= (8 << 1) { /* fine steps */
		x = (x + 1) >> 1 /* x = ceil(x / 2) */
		e++
	}
	return ((e + 1) << 3) | (x - 8)
}

// Fb2int convert back
func Fb2int(x int) int {
	if x < 8 {
		return x
	}
	return ((x & 7) + 8) << uint((x>>3)-1)
}
