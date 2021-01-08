package state

import "luago/api"

// PushNil 往Lua栈里推入nil
func (vs *LuaState) PushNil() {
	vs.stack.push(nil)
}

// PushBoolean 往Lua栈里推入布尔值
func (vs *LuaState) PushBoolean(b bool) {
	vs.stack.push(b)
}

// PushInteger 往Lua栈里推入整数
func (vs *LuaState) PushInteger(n int64) {
	vs.stack.push(n)
}

// PushNumber 往Lua栈里推入数字
func (vs *LuaState) PushNumber(n float64) {
	vs.stack.push(n)
}

// PushString 往Lua栈里推入字符串
func (vs *LuaState) PushString(s string) {
	vs.stack.push(s)
}

// PushGoFunction 往Lua栈里推入Go函数[转换为Go闭包]
func (vs *LuaState) PushGoFunction(f api.GoFunction) {
	vs.stack.push(newGoClosure(f))
}

// PushGlobalTable 将全局环境[全局表]入栈
func (vs *LuaState) PushGlobalTable() {
	global := vs.registry.get(api.LUA_RIDX_GLOBALS)
	vs.stack.push(global)

	// 通过注册表伪索引 的方式
	// vs.GetI(api.LUA_REGISTRYINDEX, api.LUA_RIDX_GLOBALS)
}
