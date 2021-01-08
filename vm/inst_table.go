package vm

import "luago/api"

const LFIELDS_PER_FLUSH = 50

// NEWTABLE指令[iABC模式]
// 创建空表，并将其放入指定寄存器，索引由A决定，初始的数组和散列表容量由B和C决定
// R(A) := {} (size = B,C)
func newTable(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

// GETTABLE指令[iABC模式]
// 根据键从表里取值，放入目标寄存器
// 表位于寄存器，索引由B指定；键可能位于寄存器也可能在常量表，由C指定； 目标寄存器由A指定
// R(A) := R(B)[RK(C)]
func getTable(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// SETTABLE指令[iABC模式]
// 根据键往表里赋值,表位于寄存器，索引由A指定，键和值位于寄存器或常量表，分别由B和C指定
// R(A)[RK(B)] := RK(C)
func setTable(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

// SETLIST指令[iABC模式 (+ EXTRAARG指令[iAx模式])
// 专门针对数组的指令，按索引批量设置数组元素
// 数组位于寄存器，索引由A指定；需要写入数组的值(紧挨数组)由B指定；数组起始索引由C决定
// R(A)[(C-1) * FPF + i] := R(A+i), 1 <= i <= B
func setList(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}
	bIsZero := b == 0
	if bIsZero {
		b = int(vm.ToInteger(-1)) - a - 1
		vm.Pop(1)
	}

	vm.CheckStack(1)
	idx := int64(c * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		vm.PushValue(a + j)
		vm.SetI(a, idx)
	}

	if bIsZero {
		for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
			idx++
			vm.PushValue(j)
			vm.SetI(a, idx)
		}

		// clear stack
		vm.SetTop(vm.RegisterCount())
	}
}
