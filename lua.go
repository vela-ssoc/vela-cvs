package csv

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
)

var xEnv vela.Environment

func cvsL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.RaiseError("invalid csv load option , must be csv.load(filename , [seek])")
		return 0
	}
	filename := L.CheckString(1)
	seek := L.IsInt(2)
	ud := newCsvGo(filename, int64(seek))
	L.Push(ud)
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env
	env.Set("csv", lua.NewFunction(cvsL))
}
