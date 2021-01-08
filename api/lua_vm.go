package api

// LuaVM lua虚拟机接口
type LuaVM interface {
	LuaState
	PC() int            // 返回当前PC(程序计数器) 仅测试使用
	AddPC(n int)        // 修改PC 用于实现跳转指令
	Fetch() uint32      // 取出当前指令，将PC指向下一条指令
	GetConst(idx int)   // 将指定常量 推入栈顶
	GetRK(rk int)       // 将指定常量或者栈值 推入栈顶
	LoadProto(idx int)  // 将当前Lua函数的子函数原型[索引由参数指定]实例化为闭包，推入栈顶
	LoadVararg(n int)   // 把传递给当前Lua函数的变长参数推入栈顶 [多退少补]
	RegisterCount() int // 返回当前Lua函数操作的寄存器数量
}
