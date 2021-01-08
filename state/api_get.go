package state

import "luago/api"

// CreateTable 创建表 估计容量
func (vs *LuaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	vs.stack.push(t)
}

// NewTable 新建表 无法估计容量
func (vs *LuaState) NewTable() {
	vs.CreateTable(0, 0)
}

// GetTable 根据键(栈顶弹出) 从表里取值，索引由参数决定，把值推入栈顶，返回值的类型
func (vs *LuaState) GetTable(idx int) api.LuaType {
	t := vs.stack.get(idx)
	k := vs.stack.pop()
	return vs.getTable(t, k)
}

func (vs *LuaState) getTable(t, k luaValue) api.LuaType {
	if tb, ok := t.(*luaTable); ok {
		v := tb.get(k)
		vs.stack.push(v)
		return typeOf(v)
	}
	panic("not a table!")
}

// GetField 根据传入的键，获取值
func (vs *LuaState) GetField(idx int, k string) api.LuaType {
	t := vs.stack.get(idx)
	return vs.getTable(t, k)

	// 利用GetTable
	// vs.PushString(k)
	// return vs.GetTable(idx)
}

// GetI 参数传入的键是数字
func (vs *LuaState) GetI(idx int, i int64) api.LuaType {
	t := vs.stack.get(idx)
	return vs.getTable(t, i)
}

// GetGlobal 将全局环境表_G.k 推入栈顶
func (vs *LuaState) GetGlobal(k string) api.LuaType {
	t := vs.registry.get(api.LUA_RIDX_GLOBALS)
	return vs.getTable(t, k)
}
