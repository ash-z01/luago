package state

import "luago/api"

// LuaState Lua机
type LuaState struct {
	registry *luaTable // 注册表
	stack    *LuaStack
}

// New 创建luaState实例
func New() *LuaState {
	registry := newLuaTable(0, 0)
	registry.put(api.LUA_RIDX_GLOBALS, newLuaTable(0, 0))
	vs := &LuaState{
		registry: registry,
	}
	vs.pushLuaStack(newLuaStack(api.LUA_MINSTACK, vs))
	return vs
}

// 入栈
func (vs *LuaState) pushLuaStack(stack *LuaStack) {
	stack.prev = vs.stack
	vs.stack = stack
}

// 出栈
func (vs *LuaState) popLuaStack() {
	stack := vs.stack
	vs.stack = stack.prev
	stack.prev = nil
}
