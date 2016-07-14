package lineFile

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	var (
		f   *File
		err error
	)

	if f, err = New(Opts{
		Path: "./",
		Name: "test",
		Ext:  "txt",
	}); err != nil {
		t.Error(err)
		return
	}

	f.WriteLine([]byte("1"))
	f.WriteLine([]byte("2"))
	f.WriteLine([]byte("3"))
	f.Flush()

	f.prevLine()
	f.prevLine()
	f.prevLine()

	f.ReadLine(func(rdr *bytes.Buffer) {
		var a [16]byte
		n, err := rdr.Read(a[:])
		fmt.Println("Read!", n, err)
		fmt.Println(string(a[:n]))
	})

	f.ReadLine(func(rdr *bytes.Buffer) {
		var a [16]byte
		n, err := rdr.Read(a[:])
		fmt.Println("Read!", n, err)
		fmt.Println(string(a[:n]))
	})

	f.ReadLine(func(rdr *bytes.Buffer) {
		var a [16]byte
		n, err := rdr.Read(a[:])
		fmt.Println("Read!", n, err)
		fmt.Println(string(a[:n]))
	})

	f.ReadLine(func(rdr *bytes.Buffer) {
		var a [16]byte
		n, err := rdr.Read(a[:])
		fmt.Println("Read!", n, err)
		fmt.Println(string(a[:n]))
	})

	f.prevLine()
	f.prevLine()
	f.ReadLine(func(rdr *bytes.Buffer) {
		var a [16]byte
		n, err := rdr.Read(a[:])
		fmt.Println("Read!", n, err)
		fmt.Println(string(a[:n]))
	})
}
