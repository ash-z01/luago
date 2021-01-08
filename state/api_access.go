package state

import (
	"fmt"
	"luago/api"
)

// TypeName ...
func (vs *LuaState) TypeName(tp api.LuaType) string {
	switch tp {
	case api.LUA_TNONE:
		return "no value"
	case api.LUA_TNIL:
		return "nil"
	case api.LUA_TBOOLEAN:
		return "boolean"
	case api.LUA_TNUMBER:
		return "number"
	case api.LUA_TSTRING:
		return "string"
	case api.LUA_TTABLE:
		return "table"
	case api.LUA_TFUNCTION:
		return "function"
	case api.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// Type ...
func (vs *LuaState) Type(idx int) api.LuaType {
	if vs.stack.isValid(idx) {
		val := vs.stack.get(idx)
		return typeOf(val)
	}
	return api.LUA_TNONE
}

func (vs *LuaState) IsNone(idx int) bool {
	return vs.Type(idx) == api.LUA_TNONE
}

func (vs *LuaState) IsNil(idx int) bool {
	return vs.Type(idx) == api.LUA_TNIL
}

func (vs *LuaState) IsNoneOrNil(idx int) bool {
	return vs.Type(idx) <= api.LUA_TNIL
}

func (vs *LuaState) IsBoolean(idx int) bool {
	return vs.Type(idx) == api.LUA_TBOOLEAN
}

func (vs *LuaState) IsString(idx int) bool {
	tp := vs.Type(idx)
	return tp == api.LUA_TSTRING || tp == api.LUA_TNUMBER
}

func (vs *LuaState) IsNumber(idx int) bool {
	_, ok := vs.ToNumberX(idx)
	return ok
}

func (vs *LuaState) IsInteger(idx int) bool {
	val := vs.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (vs *LuaState) ToBoolean(idx int) bool {
	val := vs.stack.get(idx)
	return convertToBoolean(val)
}

func (vs *LuaState) ToNumber(idx int) float64 {
	n, _ := vs.ToNumberX(idx)
	return n
}

func (vs *LuaState) ToNumberX(idx int) (float64, bool) {
	val := vs.stack.get(idx)
	return convertToFloat(val)
}

func (vs *LuaState) ToInteger(idx int) int64 {
	i, _ := vs.ToIntegerX(idx)
	return i
}

func (vs *LuaState) ToIntegerX(idx int) (int64, bool) {
	val := vs.stack.get(idx)
	return convertToInteger(val)
}

func (vs *LuaState) ToString(idx int) string {
	s, _ := vs.ToStringX(idx)
	return s
}

func (vs *LuaState) ToStringX(idx int) (string, bool) {
	val := vs.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		vs.stack.set(idx, s) // 此处会修改栈
		return s, true
	default:
		return "", false
	}
}

// IsGoFunction 判断索引处的值，是否为Go函数[闭包]
func (vs *LuaState) IsGoFunction(idx int) bool {
	val := vs.stack.get(idx)
	if c, ok := val.(*closure); ok { // 判断是否是闭包
		return c.goFunc != nil // 进一步判断是否为Go闭包
	}
	return false
}

// ToGoFunction 将索引处的值，转换为Go函数[闭包] 并返回
func (vs *LuaState) ToGoFunction(idx int) api.GoFunction {
	val := vs.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc
	}
	return nil
}
