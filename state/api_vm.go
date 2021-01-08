package state

// PC 返回当前PC
func (vs *LuaState) PC() int {
	return vs.stack.pc
}

// AddPC 修改PC (用于实现跳转指令)
func (vs *LuaState) AddPC(n int) {
	vs.stack.pc += n
}

// Fetch 取出当前指令，将PC指向下一条指令
func (vs *LuaState) Fetch() uint32 {
	i := vs.stack.closure.proto.Code[vs.stack.pc]
	vs.stack.pc++
	return i
}

// GetConst 将指定常量 推入栈顶
func (vs *LuaState) GetConst(idx int) {
	c := vs.stack.closure.proto.Constants[idx]
	vs.stack.push(c)
}

// GetRK 将指定常量或者栈值 推入栈顶
func (vs *LuaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		vs.GetConst(rk & 0xFF) // 去掉高位的1
	} else { // register
		vs.PushValue(rk + 1) // 寄存器索引从0开始，而栈索引从1开始
	}
}

// RegisterCount 返回当前Lua函数操作的寄存器数量
func (vs *LuaState) RegisterCount() int {
	return int(vs.stack.closure.proto.MaxStackSize)
}

// LoadProto 将当前Lua函数的子函数原型[索引由参数指定]实例化为闭包，推入栈顶
func (vs *LuaState) LoadProto(idx int) {
	proto := vs.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(proto)
	vs.stack.push(closure)
}

// LoadVararg 把传递给当前Lua函数的变长参数推入栈顶 [多退少补]
func (vs *LuaState) LoadVararg(n int) {
	if n < 0 {
		n = len(vs.stack.varargs)
	}
	vs.stack.check(n)
	vs.stack.pushN(vs.stack.varargs, n)
}
