package main

import (
	"luago/state"
	"testing"
)

func TestLuaState(t *testing.T) {
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

// func printStack(ls *state.LuaState) {
// 	top := ls.GetTop()
// 	for i := 1; i <= top; i++ {
// 		t := ls.Type(i)
// 		switch t {
// 		case api.LUA_TBOOLEAN:
// 			fmt.Printf("[%t]", ls.ToBoolean(i))
// 		case api.LUA_TNUMBER:
// 			fmt.Printf("[%g]", ls.ToNumber(i))
// 		case api.LUA_TSTRING:
// 			fmt.Printf("[%q]", ls.ToString(i))
// 		default:
// 			fmt.Printf("[%s]", ls.TypeName(t))
// 		}
// 	}
// 	fmt.Println()
// }
