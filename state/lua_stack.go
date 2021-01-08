package state

import "luago/api"

// LuaStack Lua栈
type LuaStack struct {
	slots   []luaValue
	top     int
	prev    *LuaStack  // 使用链表表示函数调用栈
	closure *closure   // 函数
	varargs []luaValue // 变长参数
	pc      int        // 程序计数器
	state   *LuaState  // 访问注册表
}

func newLuaStack(size int, state *LuaState) *LuaStack {
	return &LuaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: state,
	}
}

// 检查lua栈空闲空间是否还可以容纳 N个值， 如果不满足，则扩容
func (t *LuaStack) check(n int) {
	free := len(t.slots) - t.top
	for i := free; i < n; i++ {
		t.slots = append(t.slots, nil)
	}
}

// push 将值推入栈顶, 如果溢出，则panic
func (t *LuaStack) push(val luaValue) {
	if t.top == len(t.slots) {
		panic("stack overflow!")
	}
	t.slots[t.top] = val
	t.top++
}

// pop 从栈顶弹出值，如果栈是空的，则panic
func (t *LuaStack) pop() luaValue {
	if t.top < 1 {
		panic("stack underflow!")
	}
	t.top--
	val := t.slots[t.top]
	t.slots[t.top] = nil
	return val
}

// absIndex 将索引转换成 绝对索引
func (t *LuaStack) absIndex(idx int) int {
	if idx <= api.LUA_REGISTRYINDEX { // 伪索引 直接返回
		return idx
	}
	if idx >= 0 {
		return idx
	}
	return idx + t.top + 1
}

// isValid 判断 索引是否有效
func (t *LuaStack) isValid(idx int) bool {
	if idx == api.LUA_REGISTRYINDEX { // 注册表 伪索引属于有效索引
		return true
	}
	absIdx := t.absIndex(idx)
	return absIdx > 0 && absIdx <= t.top
}

// get 根据索引从栈中取值
func (t *LuaStack) get(idx int) luaValue {
	if idx == api.LUA_REGISTRYINDEX {
		// 注册表 伪索引，直接返回注册表
		return t.state.registry
	}
	absIdx := t.absIndex(idx)
	if absIdx > 0 && absIdx <= t.top {
		return t.slots[absIdx-1]
	}
	return nil
}

// set 根据索引向栈内写入值，如果索引无效 则panic
func (t *LuaStack) set(idx int, val luaValue) {
	if idx == api.LUA_REGISTRYINDEX {
		// 注册表 伪索引，直接修改注册表
		t.state.registry = val.(*luaTable)
		return
	}
	absIdx := t.absIndex(idx)
	if absIdx > 0 && absIdx <= t.top {
		t.slots[absIdx-1] = val
		return
	}
	panic("invalid index!")
}

func (t *LuaStack) reverse(from, to int) {
	slots := t.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

func (t *LuaStack) popN(n int) []luaValue {
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = t.pop()
	}
	return vals
}

func (t *LuaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}
	for i := 0; i < n; i++ {
		if i < nVals {
			t.push(vals[i])
		} else {
			t.push(nil)
		}
	}
}
