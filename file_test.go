package lineFile

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

var errInvalidValue = errors.New("invalid value")

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

	f.PrevLine()
	f.PrevLine()
	f.PrevLine()

	if err = read(f, "1"); err != nil {
		t.Error(err)
		return
	}

	if err = read(f, "2"); err != nil {
		t.Error(err)
		return
	}

	if err = read(f, "3"); err != nil {
		t.Error(err)
		return
	}

	f.PrevLine()
	f.PrevLine()
	f.PrevLine()

	if err = read(f, "1"); err != nil {
		t.Error(err)
		return
	}

	f.NextLine()

	if err = read(f, "3"); err != nil {
		t.Error(err)
		return
	}

	f.SeekToLine(0)
	if err = read(f, "1"); err != nil {
		t.Error(err)
		return
	}

	f.SeekToLine(2)
	if err = read(f, "3"); err != nil {
		t.Error(err)
		return
	}

	f.SeekToLine(1)
	if err = read(f, "2"); err != nil {
		t.Error(err)
		return
	}

	if err = f.SeekToLine(10900); err != ErrLineNotFound {
		t.Error("ErrLineNotFound is not returned when it should be")
	}

	f.Close()
	os.Remove(f.Location())
}

func read(f *File, val string) (err error) {
	f.ReadLine(func(rdr *bytes.Buffer) {
		var (
			a [16]byte
			n int
		)

		if n, err = rdr.Read(a[:]); err != nil {
			if err == io.EOF {
				err = nil
			}
		}

		if string(a[:n]) != val {
			err = errInvalidValue
		}
	})

	return
}
