package csv

import (
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"io"
)

func (c *LCvs) String() string                         { return "" }
func (c *LCvs) Type() lua.LValueType                   { return lua.LTObject }
func (c *LCvs) AssertFloat64() (float64, bool)         { return 0, false }
func (c *LCvs) AssertString() (string, bool)           { return "", false }
func (c *LCvs) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (c *LCvs) Peek() lua.LValue                       { return c }

func (c *LCvs) pipeL(L *lua.LState) int {
	r := c.Reader()
	if r == nil {
		L.Push(lua.S2L(c.err.Error()))
		return 1
	}

	pv := pipe.NewByLua(L)
	for {
		select {
		case <-c.done:
			return 0

		default:
			line, err := r.Read()
			if err != nil {
				if err == io.EOF {
					return 0
				}
				L.Push(lua.S2L(err.Error()))
				return 1
			}
			pv.Do(row(line), L, func(er error) {
				audit.Infof("%s cvs pipe call fail %v", c.filename, er).From(L.CodeVM()).Put()
			})
		}
	}
}

func (c *LCvs) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "pipe":
		return L.NewFunction(c.pipeL)

	case "next":
		return L.NewFunction(c.Next)

	case "close":
		return L.NewFunction(c.Close)

	case "err":
		return lua.S2L(c.err.Error())

	case "eof":
		return lua.LBool(c.err == io.EOF)

	default:
		return lua.LNil
	}
}
