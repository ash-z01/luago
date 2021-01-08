package state

// Len 长度 访问索引处的值，取其长度 推入栈顶
func (vs *LuaState) Len(idx int) {
	val := vs.stack.get(idx)
	if s, ok := val.(string); ok {
		vs.stack.push(int64(len(s)))
	} else if t, ok := val.(*luaTable); ok {
		vs.stack.push(int64(t.len()))
	} else {
		panic("length error!")
	}
}

// Concat 拼接 从栈顶弹出N个值，进行拼接， 推入栈顶
func (vs *LuaState) Concat(n int) {
	if n == 0 {
		vs.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if vs.IsString(-1) && vs.IsString(-2) {
				s2 := vs.ToString(-1)
				s1 := vs.ToString(-2)
				vs.stack.pop()
				vs.stack.pop()
				vs.stack.push(s1 + s2)
				continue
			}

			panic("concatenation error!")
		}
	}
	// n == 1, do nothing
}
