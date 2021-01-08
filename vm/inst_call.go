package vm

import "luago/api"

// CLOSURE指令[iABx模式]
// 将当前Lua函数的子函数原型实例化为闭包，放入A指定的寄存器
// 子函数原型来自当前函数原型的子函数原型表，索引由Bx指定
// R(A) := closure(KPROTO[Bx])
func closure(i Instruction, vm api.LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.LoadProto(bx)
	vm.Replace(a)
}

// CALL指令[iABC模式]
// 调用Lua函数，被调函数位于A指定的寄存器中
// 需要传递的参数值在紧挨被调函数的寄存器中，数量由B指定
// 调用结束后，原先存放函数和参数的寄存器，会被返回值占据，数量由C指定
// R(A), ... R(A+C-2) := R(A)(R(A+1), ... R(A+B-1))
func call(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1

	// println(":::"+ vm.StackToString())
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

// RETURN指令[iABC模式]
// 把存放在连续多个寄存器里的值返回给主调函数
// 第一个寄存器索引由A指定，寄存器数量由B指定， C未使用
// return R(A),...R(A+B-2)
func retBack(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b == 1 {
		// no return values
	} else if b > 1 {
		// b-1 return values
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		_fixStack(a, vm)
	}
}

// VARARG指令[iABC模式]
// 把传递到当前函数的变长参数加载到连续多个寄存器中
// 第一个寄存器索引由A指定，寄存器数量由B指定，C未使用
// R(A), R(A+1) ... ,R(A+B-2) = vararg
func vararg(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b != 1 { // b==0 or b>1
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

// TAILCALL指令[iABC模式]
// return R(A)(R(A+1),R(A+2),...R(A+B-1))
func tailCall(i Instruction, vm api.LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	// todo: optimize tail call!
	c := 0
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

func _pushFuncAndArgs(a, b int, vm api.LuaVM) int {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	}
	_fixStack(a, vm)
	return vm.GetTop() - vm.RegisterCount() - 1

}

func _popResults(a, b int, vm api.LuaVM) {
	if b == 1 {
		// no results
	} else if b > 1 {
		for i := a + b - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		// leave results on stack
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

func _fixStack(a int, vm api.LuaVM) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

// SELF指令[iABC模式]
// 把对象和方法拷贝到相邻的两个寄存器中
// 对象寄存器索引由B指定，方法名在常量表,索引由C指定， 目标寄存器索引由A指定
// R(A+1) := R(B); R(A) := R(B)[RK(C)]
func self(i Instruction, vm api.LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}
