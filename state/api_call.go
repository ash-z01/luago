package state

import (
	"luago/api"
	"luago/binchunk"
	"luago/vm"
)

// Load 加载二进制chunk文件, 把主函数原型实例化为闭包，推入栈顶
// Return 0成功 非0失败
// mode [b, t, bt]
func (vs *LuaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk) // todo
	c := newLuaClosure(proto)
	vs.stack.push(c)
	return 0 // TODO
}

// Call 调用Lua函数，执行之前必须先把被调函数推入栈顶，然后把参数值依次推入栈顶
// 调用结束后，参数值和函数会被弹出栈顶，取而代之的是指定数量的返回值
// nArgs: 传递给被调函数的参数数量，[同时也隐含给出了被调函数在栈里的位置]
// nRets: 指定了需要的返回值数量[多退少补]，如果是-1, 则被调函数的返回值全部留在栈顶
func (vs *LuaState) Call(nArgs, nRets int) {
	val := vs.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		// fmt.Printf("call %s<%d,%d>\n", c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)
		if c.proto != nil {
			vs.callLuaClosure(nArgs, nRets, c)
		} else {
			// 调用Go函数[闭包]
			vs.callGoClosure(nArgs, nRets, c)
		}
	} else {
		panic("not function!")
	}
}

func (vs *LuaState) callLuaClosure(nArgs, nRets int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+api.LUA_MINSTACK, vs)
	newStack.closure = c

	// pass args, pop func
	funcAndArgs := vs.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// run closure
	vs.pushLuaStack(newStack)
	vs.runLuaClosure()
	vs.popLuaStack()

	// return results
	if nRets != 0 {
		results := newStack.popN(newStack.top - nRegs)
		vs.stack.check(len(results))
		vs.stack.pushN(results, nRets)
	}
}

func (vs *LuaState) runLuaClosure() {
	for {
		inst := vm.Instruction(vs.Fetch())
		inst.Execute(vs)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

func (vs *LuaState) callGoClosure(nArgs, nRets int, c *closure) {
	// create new lua stack
	newStack := newLuaStack(nArgs+api.LUA_MINSTACK, vs)
	newStack.closure = c

	// pass args, pop func
	if nArgs > 0 {
		args := vs.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	vs.stack.pop()

	// run closure
	vs.pushLuaStack(newStack)
	r := c.goFunc(vs)
	vs.popLuaStack()

	// return results
	if nRets != 0 {
		results := newStack.popN(r)
		vs.stack.check(len(results))
		vs.stack.pushN(results, nRets)
	}
}
