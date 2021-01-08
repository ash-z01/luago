package binchunk

import (
	"encoding/binary"
	"math"
)

// 二进制chunk解析结构体
type reader struct {
	data []byte
}

// readByte 读取一个字节 1byte
func (r *reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

// readUint32 以小端方式 读取cint存储类型 4byte
func (r *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

// readUint64 以小端方式 读取size_t存储类型 8byte
func (r *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

// readLuaInteger 读取一个Lua整数 8byte
func (r *reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

// readLuaNumber 读取一个Lua浮点数 8byte
func (r *reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

// readString 读取字符串
func (r *reader) readString() string {
	size := uint(r.readByte()) // 短字符串?
	if size == 0 {             // NULL字符
		return ""
	}
	if size == 0xFF { // 长字符串
		size = uint(r.readUint64())
	}
	bytes := r.readBytes(size - 1)
	return string(bytes)
}

// readBytes 读取n个字节
func (r *reader) readBytes(n uint) []byte {
	bytes := r.data[:n]
	r.data = r.data[n:]
	return bytes
}

// checkHeader 检查头部
func (r *reader) checkHeader() {
	if string(r.readBytes(4)) != LuaSignature {
		panic("not a precompiled chunk!")
	}
	if r.readByte() != LuacVersion {
		panic("version misnatch!")
	}
	if r.readByte() != LuacFormat {
		panic("format mismatch!")
	}
	if string(r.readBytes(6)) != LuacData {
		panic("corrupted!")
	}
	if r.readByte() != CintSize {
		panic("int size mismatch!")
	}
	if r.readByte() != CsizetSize {
		panic("size_t size mismatch!")
	}
	if r.readByte() != InstructionSize {
		panic("instruction size mismatch!")
	}
	if r.readByte() != LuaIntegerSize {
		panic("lua_Integer size mismatch!")
	}
	if r.readByte() != LuaNumberSize {
		panic("lua_Number size mismatch!")
	}
	if r.readLuaInteger() != LuacInt {
		panic("endianness mismatch!")
	}
	if r.readLuaNumber() != LuacNum {
		panic("float format mismatch!")
	}
}

// readProto 读取函数原型
func (r *reader) readProto(parentSource string) *Prototype {
	source := r.readString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     r.readUint32(),
		LastLineDefined: r.readUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStackSize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}

// readCode 读取指令表
func (r *reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}
	return code
}

// readConstants 读取常量表
func (r *reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

// 从字节流里读取一个常量
func (r *reader) readConstant() interface{} {
	switch r.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return r.readByte() != 0
	case TAG_INTEGER:
		return r.readLuaInteger()
	case TAG_NUMBER:
		return r.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR:
		return r.readString()
	default:
		panic("corrupted!") // todo
	}
}

// readUpvalues 从字节流里读取Upvalue表
func (r *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return upvalues
}

// readProtos 从字节流里读取函数原型表 函数原型->递归结构
func (r *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, r.readUint32())
	for i := range protos {
		protos[i] = r.readProto(parentSource)
	}
	return protos
}

// readLineInfo 从字节流 读取行号表
func (r *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, r.readUint32())
	for i := range lineInfo {
		lineInfo[i] = r.readUint32()
	}
	return lineInfo
}

// readLocVars 从字节流 读取局部变量表
func (r *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, r.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}
	return locVars
}

// readUpvalueNames 从字节流 读取Upvalue名列表
func (r *reader) readUpvalueNames() []string {
	names := make([]string, r.readUint32())
	for i := range names {
		names[i] = r.readString()
	}
	return names
}
