package api

/* basic types */
const (
	LUA_TNONE = iota - 1 // -1; Lua栈可以按索引取值，如果是无效值，返回 LUA_TNONE
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

/* arithmetic functions */
const (
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // *
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // //
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ~
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPUNM         // -
	LUA_OPBNOT        // ~
)

/* comparison functions */
const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        // <=
)

// LuaState
const (
	LUA_MINSTACK            = 20
	LUAI_MAXSTACK           = 1000000
	LUA_REGISTRYINDEX       = -LUAI_MAXSTACK - 1000
	LUA_RIDX_GLOBALS  int64 = 2
)
