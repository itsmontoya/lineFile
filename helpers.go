package lineFile

import (
	"fmt"
	"os"
)

// Opts are used to initialize a file
type Opts struct {
	Path string
	Name string
	Ext  string

	NoSet bool // If set to true, will avoid setting file when calling New
}

func (o Opts) isValid() bool {
	return len(o.Path) > 0 && len(o.Name) > 0 && len(o.Ext) > 0
}

func peek(f *os.File) {
	var pkk [32]byte
	n, _ := f.Read(pkk[:])
	f.Seek(int64(-n), 1)
	fmt.Println("Peeek??", pkk[:n])
}

func reverseByteSlice(bs []byte) {
	var n int
	c := len(bs) - 1
	for i := range bs {
		if n = c - i; n == i || n < i {
			break
		}

		bs[i], bs[n] = bs[n], bs[i]
	}
}

func getNewlineIndex(s []byte) (idx int) {
	for i, b := range s {
		if b == charNewline {
			return i
		}
	}

	return -1
}
