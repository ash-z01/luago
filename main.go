package main

import (
	"fmt"
	"io/ioutil"
	"luago/api"
	"luago/binchunk"
	"luago/state"
	"luago/vm"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		// proto := binchunk.Undump(data)
		// list(proto)
		// luaMain(proto)
		ls := state.New()
		ls.Register("print", print)
		ls.Load(data, os.Args[1], "b")
		ls.Call(0, 0)
	}
}

func print(vs api.LuaState) int {
	nArgs := vs.GetTop()
	for i := 1; i <= nArgs; i++ {
		if vs.IsBoolean(i) {
			fmt.Printf("%t", vs.ToBoolean(i))
		} else if vs.IsString(i) {
			fmt.Printf("%s", vs.ToString(i))
		} else {
			fmt.Print(vs.TypeName(vs.Type(i)))
		}
		if i < nArgs {
			fmt.Printf("\t")
		}
	}
	fmt.Println()
	return 0
}

// func luaMain(proto *binchunk.Prototype) {
// 	nRegs := int(proto.MaxStackSize)
// 	ls := state.New(nRegs+8, proto)
// 	ls.SetTop(nRegs)
// 	for {
// 		pc := ls.PC()
// 		inst := vm.Instruction(ls.Fetch())
// 		if inst.Opcode() != vm.OP_RETURN {
// 			inst.Execute(ls)
// 			fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
// 			printStack(ls)
// 		} else {
// 			break
// 		}
// 	}
// }

func printStack(ls api.LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case api.LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case api.LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func list(f *binchunk.Prototype) {
	printHeader(f)
	printCode(f)
	printDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}

func printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}
	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}
	fmt.Printf("\n%s <%s:%d, %d> (%d instructions)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))
	fmt.Printf("%d%s params, %d slots, %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))
	fmt.Printf("%d locals, %d constants, %d functions\n",
		len(f.LocVars), len(f.Constants), len(f.Protos))
}

func printCode(f *binchunk.Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		i := vm.Instruction(c)
		// fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, c)
		fmt.Printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Println()
	}
}

func printOperands(i vm.Instruction) {
	switch i.OpMode() {
	case vm.IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != vm.OpArgN {
			if b > 0XFF {
				fmt.Printf(" %d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != vm.OpArgN {
			if c > 0XFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case vm.IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		if i.BMode() == vm.OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == vm.OpArgU {
			fmt.Printf(" %d", bx)
		}
	case vm.IAsBx:
		a, sBx := i.AsBx()
		fmt.Printf("%d %d", a, sBx)
	case vm.IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)
	default:
	}
}

func printDetail(f *binchunk.Prototype) {
	fmt.Printf("constants (%d)\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}
	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}
	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalName(f, i), upval.Instack, upval.Idx)
	}
}

// 将常量表中的常量 转换为 字符串
func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

// 根据Upvalue索引从调试信息里找出Upvalue的名字
func upvalName(f *binchunk.Prototype, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}

// luaState lua栈测试
func testLuaState1() {
	ls := state.New()
	ls.PushBoolean(true)
	printStack(ls)
	ls.PushInteger(10)
	printStack(ls)
	ls.PushNil()
	printStack(ls)
	ls.PushString("hello")
	printStack(ls)
	ls.PushValue(-4)
	printStack(ls)
	ls.Replace(3)
	printStack(ls)
	ls.SetTop(6)
	printStack(ls)
	ls.Remove(-3)
	printStack(ls)
	ls.SetTop(-5)
	printStack(ls)
}

// LuaState 运算测试
func testLuaState2() {
	ls := state.New()
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	printStack(ls)

	ls.Arith(api.LUA_OPADD)
	printStack(ls)
	ls.Arith(api.LUA_OPBNOT)
	printStack(ls)
	ls.Len(2)
	printStack(ls)
	ls.Concat(3)
	printStack(ls)
	ls.PushBoolean(ls.Compare(1, 2, api.LUA_OPEQ))
	printStack(ls)
}
