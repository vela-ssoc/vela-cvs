package csv

import (
	"encoding/csv"
	"github.com/vela-ssoc/vela-kit/lua"
	"io"
	"os"
)

type LCvs struct {
	done chan struct{}

	filename string
	fd       *os.File
	reader   *csv.Reader
	err      error
	eof      bool
	seek     int64
}

func newCsvGo(filename string, seek int64) *LCvs {
	return &LCvs{
		done:     make(chan struct{}, 1),
		filename: filename,
		seek:     seek,
	}
}

func (c *LCvs) Fd() *os.File {
	if c.fd != nil {
		return c.fd
	}

	fd, err := os.Open(c.filename)
	if err != nil {
		c.err = err
		return fd
	}

	if c.seek != 0 {
		if _, e := fd.Seek(c.seek, io.SeekStart); e != nil {
			c.err = e
			return nil
		}
	}

	c.fd = fd
	return fd
}

func (c *LCvs) Reader() *csv.Reader {
	if c.reader != nil {
		return c.reader
	}

	fd := c.Fd()
	if fd == nil {
		return nil
	}

	reader := csv.NewReader(fd)
	c.reader = reader

	return reader
}

func (c *LCvs) Next(L *lua.LState) int {
	r := c.Reader()
	if r == nil {
		L.Push(lua.LNil)
		L.Push(lua.S2L(c.err.Error()))
		return 2
	}

	line, e := r.Read()
	if e != nil {
		c.err = e
		L.Push(lua.LNil)
		L.Push(lua.S2L(e.Error()))
		return 2
	}
	L.Push(row(line))
	return 1
}

func (c *LCvs) Close(L *lua.LState) int {
	c.done <- struct{}{}

	if c.fd == nil {
		return 0
	}

	if e := c.fd.Close(); e != nil {
		xEnv.Errorf("cvs %s file close error %v", c.filename, e)
		return 0
	}

	xEnv.Errorf("cvs %s file close succeed", c.filename)
	return 0
}
