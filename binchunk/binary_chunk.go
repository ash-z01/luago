package binchunk

const (
	// LuaSignature LUA_SIGNATURE    = "\x1bLua"
	LuaSignature = "\x1bLua"
	// LuacVersion LUAC_VERSION     = 0x53
	LuacVersion = 0x53
	// LuacFormat LUAC_FORMAT      = 0
	LuacFormat = 0
	// LuacData LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	LuacData = "\x19\x93\r\n\x1a\n"
	// CintSize CINT_SIZE        = 4
	CintSize = 4
	// CsizetSize CSIZET_SIZE      = 8
	CsizetSize = 8
	// InstructionSize INSTRUCTION_SIZE = 4
	InstructionSize = 4
	// LuaIntegerSize LUA_INTEGER_SIZE = 8
	LuaIntegerSize = 8
	// LuaNumberSize LUA_NUMBER_SIZE  = 8
	LuaNumberSize = 8
	// LuacInt LUAC_INT         = 0X5678
	LuacInt = 0x5678
	// LuacNum LUUAC_NUM        = 370.5
	LuacNum = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type binaryChunk struct {
	header                  // 头部
	sizeUpvalues byte       // 主函数upvalue数量
	mainFunc     *Prototype // 主函数原型
}

type header struct {
	// lua二进制格式开头的魔数/签名
	// 分别是 ESC, L, u, a 的ASCII码
	// 0x1B4C7561  ==  Go的字面量 "\x1bLua"
	signature [4]byte

	// 版本号 大版本MajorVer 小版本MinorVer 发布版本ReleaseVer
	// 5.3.4  version = majorVer x 16 + minorVer = 5 x 16 + 3 = 83 = 0x53
	version byte

	// 格式号 检查是否匹配，官方使用的是0
	format byte

	// 0x1993 Lua1.0发布的年份  回车符[0x0D] 换行符[0x0A] 替换符[0x1A] 换行符[0x0A]
	// Go字面量 "\x19\x93\r\n\x1a\n"
	// 用来做进一步的校验
	luacData [6]byte

	// cint, size_t, Lua虚拟机指令，Lua整型 Lua浮点型 分别占用的字节
	// 0x04   0x08    0x04         0x08     0x08
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte

	// 接下来的n个字节存放Lua整数值0x5678
	// 检测大小端方式
	luacInt int64

	// 头部的最后n个字节存放Lua浮点数370.5
	luacNum float64
}

// Upvalue ...
type Upvalue struct {
	Instack byte
	Idx     byte
}

// LocVar 局部变量
type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

// Prototype ... 函数原型
type Prototype struct {
	// 源文件名
	Source string
	// 起始行号
	LineDefined uint32
	// 终止行号
	LastLineDefined uint32
	// 固定参数的个数
	NumParams byte
	// 是否是vararg函数
	IsVararg byte
	// 寄存器数量
	MaxStackSize byte
	// 指令表
	Code []uint32
	// 常量表
	Constants []interface{}
	// Upvalue表
	Upvalues []Upvalue
	// 子函数原型表
	Protos []*Prototype
	// 行号表
	LineInfo []uint32
	// 局部变量表
	LocVars []LocVar
	// Upvalue名 列表
	UpvalueNames []string
}

// Undump 解析二进制 chunk
func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        // 检查头部
	reader.readByte()           // 跳过Upvalue数量
	return reader.readProto("") // 读取函数原型
}
