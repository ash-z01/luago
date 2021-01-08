package vm

import "luago/api"

const MAXARG_Bx = 1<<18 - 1       // 2^18 - 1 = 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 262143 / 2 = 131071

// Instruction 指令类型
/*
 31       22       13       5    0
  +-------+^------+-^-----+-^-----
  |b=9bits |c=9bits |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    bx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |   sbx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    ax=26bits            |op=6|
  +-------+^------+-^-----+-^-----
 31      23      15       7      0
*/
type Instruction uint32

// Opcode 提取操作码
func (i Instruction) Opcode() int {
	return int(i & 0x3F)
}

// ABC 从iABC指令模式中提取参数
func (i Instruction) ABC() (a, b, c int) {
	a = int(i >> 6 & 0xFF)
	c = int(i >> 14 & 0x1FF)
	b = int(i >> 23 & 0x1FF)
	return
}

// ABx 从iABx指令模式中提取参数
func (i Instruction) ABx() (a, bx int) {
	a = int(i >> 6 & 0xFF)
	bx = int(i >> 14)
	return
}

// AsBx 从iAsBx指令模式中提取参数
func (i Instruction) AsBx() (a, sBx int) {
	a, bx := i.ABx()
	return a, bx - MAXARG_sBx
}

// Ax 从iAx指令模式中提取参数
func (i Instruction) Ax() int {
	return int(i >> 6)
}

// OpName 读取指令操作码名字
func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

// OpMode 读取指令编码模式
func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}

// BMode 读取指令操作数B的使用模式
func (i Instruction) BMode() byte {
	return opcodes[i.Opcode()].argBMode
}

// CMode 读取指令操作数C的使用模式
func (i Instruction) CMode() byte {
	return opcodes[i.Opcode()].argCMode
}

// Execute 指令执行
func (i Instruction) Execute(vm api.LuaVM) {
	action := opcodes[i.Opcode()].action
	if action != nil {
		action(i, vm)
	} else {
		panic(i.OpName())
	}
}
