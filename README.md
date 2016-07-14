# Line file [![GoDoc](https://godoc.org/github.com/itsmontoya/hippy?status.svg)](https://godoc.org/github.com/itsmontoya/hippy) ![Status](https://img.shields.io/badge/status-alpha-red.svg)
Line file is a file management assistant for line-based operations

## Usage
```go
package main

import (
	"bytes"
	"fmt"
)

func main() {
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
```