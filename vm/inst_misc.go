package vm

import . "luago/api"

// MOVE指令[iABC模式] 把源寄存器的值 移动(复制) 到目标寄存器
// R(A) := R(B)
func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1 // 寄存器索引从0开始， 栈索引从1开始
	b += 1
	vm.Copy(b, a)
}

// JMP指令[iAsBx模式] TODO
func jmp(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		panic("todo!")
	}
}
