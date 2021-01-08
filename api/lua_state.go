package api

// LuaType 类型
type LuaType = int

// ArithOp 算术与位运算相关类型
type ArithOp = int

// CompareOp 比较运算相关类型
type CompareOp = int

// GoFunction 与Go函数交互
type GoFunction = func(LuaState) int

// LuaState (Lua <-> Go 交互 接口)
type LuaState interface {
	// basic statck manipulation 基础栈操作
	GetTop() int
	AbsIndex(idx int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromIdx, toIdx int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx, n int)
	SetTop(idx int)

	// access function (stack -> Go)访问函数
	TypeName(tp LuaType) string
	Type(idx int) LuaType
	IsNone(idx int) bool
	IsNil(idx int) bool
	IsNoneOrNil(idx int) bool
	IsBoolean(idx int) bool
	IsInteger(idx int) bool
	IsNumber(idx int) bool
	IsString(idx int) bool
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)

	// push function (Go -> Stack)
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)
	/* Comparison and arithmetic functions */
	Arith(op ArithOp)
	Compare(idx1, idx2 int, op CompareOp) bool
	/* get functions (Lua -> stack) */
	NewTable()
	CreateTable(nArr, nRec int)
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType
	/* set functions (stack -> Lua) */
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, i int64)
	/* miscellaneous functions */
	Len(idx int)
	Concat(n int)

	/* load & call */

	// 加载二进制chunk文件, 把主函数原型实例化为闭包，推入栈顶
	// 实际上不仅仅可以加载 二进制chunk，也可以加载lua脚本
	// 如果无法加载chunk，需要在栈顶留下错误信息
	// return 0 成功 !0 失败
	Load(chunk []byte, chunkName, mode string) int

	// 调用Lua函数，执行之前必须先把被调函数推入栈顶，然后把参数值依次推入栈顶
	// Call结束后，参数值和函数会被弹出栈顶，取而代之的是指定数量的返回值
	// nArgs: 传递给被调函数的参数数量，[同时也隐含给出了被调函数在栈里的位置]
	// nRets: 指定了需要的返回值数量[多退少补]，如果是-1, 则被调函数的返回值全部留在栈顶
	Call(nArgs, nRets int)

	// GoFunction <---> LuaState
	PushGoFunction(f GoFunction)
	IsGoFunction(idx int) bool
	ToGoFunction(idx int) GoFunction

	// Global env
	PushGlobalTable()
	GetGlobal(name string) LuaType
	SetGlobal(name string)
	Register(name string, f GoFunction)
}
