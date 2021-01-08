package vm

import . "luago/api"

// FOR循环指令
// 数值for 和 通用for
// 数值for 按照一定步长，遍历某个范围内的数值
// 通用for 主要用于遍历表

// 数值FOR
// FORPREP 和 FORLOOP

// FORPREP指令[iAsBx模式]
// R(A) -= R(A+2);
// PC += sBx

// FORLOOP指令[iAsBx模式]
// R(A) += R(A+2)
// if R(A) <?= R(A+1) then {
//     PC += sBx;
//     R(A+3) = R(A)
// }
// <?= 当步长是正数,则为 <=，步长负数 则为 >=

/*
// -----------------------------------------------------------------
luac -l -l -
for i=1,100,2 do f() end
^D [Ctrl+D]
// -----------------------------------------------------------------
main <stdin:0,0> (8 instructions at 0x558743a85e20)
0+ params, 5 slots, 1 upvalue, 4 locals, 4 constants, 0 functions
        1       [1]     LOADK           0 -1    ; 1
        2       [1]     LOADK           1 -2    ; 100
        3       [1]     LOADK           2 -3    ; 2
        4       [1]     FORPREP         0 2     ; to 7
        5       [1]     GETTABUP        4 0 -4  ; _ENV "f"
        6       [1]     CALL            4 1 1
        7       [1]     FORLOOP         0 -3    ; to 5
        8       [1]     RETURN          0 1
constants (4) for 0x558743a85e20:
        1       1
        2       100
        3       2
        4       "f"
locals (4) for 0x558743a85e20:
        0       (for index)     4       8
        1       (for limit)     4       8
        2       (for step)      4       8
        3       i       5       7
upvalues (1) for 0x558743a85e20:
		0       _ENV    1       0
// -----------------------------------------------------------------
*/

// locals中三个特殊的局部变量 (for index) (for limit) (for step)
// 分别对应 R(A)， R(A+1)， R(A+2)
// 临时变量i 对应 R(A+3)
// ∴ FORPREP指令 在循环开始之前，预先给数值减去步长，然后跳转到FORLOOP 进入循环
// FORLOOP指令则是给数值加上步长，然后判断是否在范围之内
// 如果超过范围，则终止循环；否则把数值复制到用户定义的局部变量，然后进入循环体执行内部代码

// FORPREP指令[iAsBx模式]
// R(A) -= R(A+2);
// PC += sBx
func forPrep(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1
	// todo 可能对 a 和 b 类型转换
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(a)
	vm.AddPC(sBx)
}

// FORLOOP指令[iAsBx模式]
// R(A) += R(A+2)
// if R(A) <?= R(A+1) then {
//     PC += sBx;
//     R(A+3) = R(A)
// }
// <?= 当步长是正数,则为 <=，步长负数 则为 >=
func forLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)
	vm.Replace(a)

	isPositiveStep := vm.ToNumber(a+2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {
		// pc+=sBx; R(A+3)=R(A)
		vm.AddPC(sBx)
		vm.Copy(a, a+3)
	}
}
