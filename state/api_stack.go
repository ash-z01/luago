package state

func (vs *LuaState) GetTop() int {
	return vs.stack.top
}

// AbsIndex ...
func (vs *LuaState) AbsIndex(idx int) int {
	return vs.stack.absIndex(idx)
}

func (vs *LuaState) CheckStack(n int) bool {
	vs.stack.check(n)
	return true
}

func (vs *LuaState) Pop(n int) {
	for i := 0; i < n; i++ {
		vs.stack.pop()
	}
	// vs.SetTop(-n - 1)
}

func (vs *LuaState) Copy(fromIdx, toIdx int) {
	val := vs.stack.get(fromIdx)
	vs.stack.set(toIdx, val)
}

func (vs *LuaState) PushValue(idx int) {
	val := vs.stack.get(idx)
	vs.stack.push(val)
}

func (vs *LuaState) Replace(idx int) {
	val := vs.stack.pop()
	vs.stack.set(idx, val)
}

func (vs *LuaState) Insert(idx int) {
	vs.Rotate(idx, 1)
}

func (vs *LuaState) Remove(idx int) {
	vs.Rotate(idx, -1)
	vs.Pop(1)
}

// Rotate n 方向 >0 朝向栈顶 <=0 朝栈底
func (vs *LuaState) Rotate(idx int, n int) {
	t := vs.stack.top - 1           /* end of stack segment being rotated */
	p := vs.stack.absIndex(idx) - 1 /* start of segment */
	var m int                       /* end of prefix */
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	vs.stack.reverse(p, m)   /* reverse the prefix with length 'n' */
	vs.stack.reverse(m+1, t) /* reverse the suffix */
	vs.stack.reverse(p, t)   /* reverse the entire segment */
}

func (vs *LuaState) SetTop(idx int) {
	newTop := vs.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}
	n := vs.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			vs.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			vs.stack.push(nil)
		}
	}
}
