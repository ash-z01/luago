package vm

import . "luago/api"

// LOADNIL指令[iABC模式] 给连续N个寄存器设置 nil值，起始索引由A指定，数量由B指定，C未使用
// R(A), R(A+1), R(A+2) ... R(A+B) := nil
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

// LOADBOOLEAN指令[iABC模式]
// 给单个寄存器设置布尔值，索引由A指定，布尔数由B指定[0=false !0=ture] 如果寄存器C非0，则跳过下一条指令
// R(A) := (bool)B; if (C) then PC++
func loadBoolean(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}
}

// LOADK指令[iABx模式]
// 将常量表中某常量，加载到指定寄存器，寄存器索引由A指定，常量表索引由Bx指定
// R(A) := Kst(Bx)  // Kst 常量表
func loadK(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1
	vm.GetConst(bx)
	vm.Replace(a)
}

// LOADKX指令[iABx模式] + EXTRAARG指令[iAx模式]
func loadKX(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a += 1
	ax := Instruction(vm.Fetch()).Ax()

	//vm.CheckStack(1)
	vm.GetConst(ax)
	vm.Replace(a)
}
