package state

import (
	"luago/number"
	"math"
)

type luaTable struct {
	arr  []luaValue
	dict map[luaValue]luaValue
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t.dict = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (t *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(t.arr)) {
			return t.arr[idx-1]
		}
	}
	return t.dict[key]
}

func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (t *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(t.arr))
		if idx <= arrLen {
			t.arr[idx-1] = val
			if idx == arrLen && val == nil {
				t._shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(t.dict, key)
			if val != nil {
				t.arr = append(t.arr, val)
				t._expandArray()
			}
			return
		}
	}
	if val != nil {
		if t.dict == nil {
			t.dict = make(map[luaValue]luaValue, 8)
		}
		t.dict[key] = val
	} else {
		delete(t.dict, key)
	}
}

func (t *luaTable) _shrinkArray() {
	for i := len(t.arr) - 1; i >= 0; i-- {
		if t.arr[i] == nil {
			t.arr = t.arr[0:i]
		}
	}
}

func (t *luaTable) _expandArray() {
	for idx := int64(len(t.arr)) + 1; true; idx++ {
		if val, found := t.dict[idx]; found {
			delete(t.dict, idx)
			t.arr = append(t.arr, val)
		} else {
			break
		}
	}
}

func (t *luaTable) len() int {
	return len(t.arr)
}
