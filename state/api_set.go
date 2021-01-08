package state

import "luago/api"

// SetTable 将键值对写入表，[键和值从栈里弹出，表则由索引指定]
func (vs *LuaState) SetTable(idx int) {
	t := vs.stack.get(idx)
	v := vs.stack.pop()
	k := vs.stack.pop()
	vs.setTable(t, k, v)
}

func (vs *LuaState) setTable(t, k, v luaValue) {
	if tb, ok := t.(*luaTable); ok {
		tb.put(k, v)
		return
	}
	panic("not a table!")
}

// SetField 同SetTable, 但键改为由参数传入的字符串
func (vs *LuaState) SetField(idx int, k string) {
	t := vs.stack.get(idx)
	v := vs.stack.pop()
	vs.setTable(t, k, v)
}

// SetI 同SetTable, 但键改为由参数传入的数字
func (vs *LuaState) SetI(idx int, i int64) {
	t := vs.stack.get(idx)
	v := vs.stack.pop()
	vs.setTable(t, i, v)
}

// SetGlobal 将栈顶值弹出 写入 全局环境表_G.k
func (vs *LuaState) SetGlobal(k string) {
	t := vs.registry.get(api.LUA_RIDX_GLOBALS)
	v := vs.stack.pop()
	vs.setTable(t, k, v)
}

// Register 用于给全局环境注册Go函数
func (vs *LuaState) Register(name string, f api.GoFunction) {
	vs.PushGoFunction(f)
	vs.SetGlobal(name)
}
