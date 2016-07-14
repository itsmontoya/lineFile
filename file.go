package lineFile

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/missionMeteora/toolkit/bufferPool"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// Buffer size used for file seeking
	seekerBufSize = 32
	// Default size of byte.Buffers produced by bufferPool
	bufferPoolSize = 32
	// Byte for newline
	charNewline = '\n'
)

const (
	// ErrLineNotFound is returned when a line is not found while calling SeekNextLine
	ErrLineNotFound = errors.Error("line not found")

	// ErrIsClosed is returned when an action is attempted on a closed instance
	ErrIsClosed = errors.Error("cannot perform action on closed instance")

	// ErrIsOpen is returned when an instance is attempted to be re-opened when it's already active
	ErrIsOpen = errors.Error("cannot open an instance which is already open")

	// ErrInvalidOptions is returned when options are invalid
	ErrInvalidOptions = errors.Error("options are invalid")
)

var (
	// LineFile-global buffer pool
	bp = bufferPool.New(bufferPoolSize)
)

// New will return a pointer to a new instance of File
func New(o Opts) (f *File, err error) {
	if !o.isValid() {
		err = ErrInvalidOptions
		return
	}

	// Assign f with pointer to a new basic File
	f = &File{
		fLoc:   filepath.Join(o.Path, o.Name+"."+o.Ext),
		closed: true, // This will be set to false by f.Open
	}

	// If NoSet option is set to true, return early
	if o.NoSet {
		return
	}

	err = f.Open()
	return
}

// File is a line-based file for easy management
type File struct {
	mux sync.Mutex

	// File location (path, name, and extension )
	fLoc string

	// Seek buffer, used for storing read data while seeking
	seekBuf [seekerBufSize]byte

	// File
	f *os.File
	// Write buffer
	buf *bufio.Writer

	// Closed state, file can be re-opened using f.Open
	closed bool
}

func (f *File) getPosition() (pos int64) {
	pos, _ = f.f.Seek(0, os.SEEK_CUR)
	return
}

func (f *File) seekBackwards(cc int64) (nc int64, err error) {
	if cc > seekerBufSize {
		cc = seekerBufSize
	}

	return f.f.Seek(-cc, os.SEEK_CUR)
}

func (f *File) readChunks(fn func(int) bool) (err error) {
	var n int

	for n, err = f.f.Read(f.seekBuf[:]); ; n, err = f.f.Read(f.seekBuf[:]) {
		if err == io.EOF && n == 0 {
			err = nil
			break
		} else if err != nil {
			break
		}

		if fn(n) {
			break
		}
	}

	return
}

func (f *File) readReverseChunks(fn func(int) bool) (err error) {
	var (
		curr = f.getPosition()

		n    int
		cc   int64
		done bool
	)

	for !done {
		if curr > seekerBufSize {
			cc = seekerBufSize
		} else {
			cc = curr
			done = true
		}

		if curr, err = f.seekBackwards(curr); err != nil {
			return
		}

		if n, err = f.f.Read(f.seekBuf[:cc]); err == io.EOF && n == 0 {
			err = nil
			break
		} else if err != nil {
			break
		}

		if fn(n) {
			break
		}

		if done {
			break
		}

		if curr, err = f.seekBackwards(curr); err != nil {
			break
		}
	}

	return
}

func (f *File) nextLine() (err error) {
	var (
		nlf    bool
		offset int64 = -1
	)

	pcfn := func(n int) (end bool) {
		for i, b := range f.seekBuf[:n] {
			if b == charNewline {
				nlf = true
			} else if nlf {
				offset = int64(n - i)
				return true
			}
		}

		return
	}

	if err = f.readChunks(pcfn); err != nil {
		return
	}

	if offset == -1 {
		return ErrLineNotFound
	}

	_, err = f.f.Seek(-offset, os.SEEK_CUR)
	return
}

func (f *File) prevLine() (err error) {
	var (
		nlc    int
		offset int64 = -1
	)

	pcfn := func(n int) (end bool) {
		s := f.seekBuf[:n]
		reverseByteSlice(s)

		for i, b := range s {
			if b != charNewline {
				continue
			}

			if nlc++; nlc == 2 {
				offset = int64(i)
				return true
			}
		}

		return
	}

	if err = f.readReverseChunks(pcfn); err != nil {
		return
	}

	if offset == -1 {
		_, err = f.f.Seek(0, os.SEEK_SET)
	} else {
		_, err = f.f.Seek(-offset, os.SEEK_CUR)
	}

	return
}

func (f *File) readLine(fn func(*bytes.Buffer)) (err error) {
	var (
		n   int
		s   []byte
		idx int
		buf = bp.Get()
	)

	for err == nil {
		if n, err = f.f.Read(f.seekBuf[:]); err != nil && err != io.EOF {
			break
		}

		s = f.seekBuf[:n]
		if idx = getNewlineIndex(s); idx == -1 {
			buf.Write(s)
			continue
		}

		buf.Write(s[:idx])
		f.f.Seek(int64(-(n - idx - 1)), os.SEEK_CUR)
		err = nil
		fn(buf)
		break
	}

	bp.Put(buf)
	return
}

// Open will open a closed File
func (f *File) Open() (err error) {
	f.mux.Lock()
	if !f.closed {
		err = ErrIsOpen
		goto END
	}

	// Open persistance file
	if f.f, err = os.OpenFile(f.fLoc, os.O_CREATE|os.O_RDWR, 0644); err != nil {
		goto END
	}

	f.buf = bufio.NewWriter(f.f)
	f.closed = false

END:
	f.mux.Unlock()
	return
}

// SeekToStart will seek the file to the start
func (f *File) SeekToStart() (err error) {
	f.mux.Lock()
	if !f.closed {
		err = ErrIsOpen
		goto END
	}

	_, err = f.f.Seek(0, os.SEEK_SET)

END:
	f.mux.Unlock()
	return
}

// SeekToEnd will seek the file to the end
func (f *File) SeekToEnd() (err error) {
	f.mux.Lock()
	if !f.closed {
		err = ErrIsOpen
		goto END
	}

	_, err = f.f.Seek(0, os.SEEK_END)

END:
	f.mux.Unlock()
	return
}

// WriteLine will write a line given a provided body
func (f *File) WriteLine(b []byte) (err error) {
	f.mux.Lock()
	if f.closed {
		err = ErrIsClosed
		goto END
	}

	if _, err = f.buf.Write(b); err != nil {
		goto END
	}

	// Write our suffix byte (newline) without any middlewares so we can find a line-end without decoding
	err = f.buf.WriteByte(charNewline)

END:
	f.mux.Unlock()
	return
}

// Flush will flush the internal buffer and sync the file to disk
func (f *File) Flush() (err error) {
	f.mux.Lock()
	if f.closed {
		err = ErrIsClosed
		goto END
	}

	if err = f.buf.Flush(); err != nil {
		return
	}

	err = f.f.Sync()

END:
	f.mux.Unlock()
	return
}

// ReadLine will return a line in the form of an a bytes.Buffer
func (f *File) ReadLine(fn func(*bytes.Buffer)) (err error) {
	f.mux.Lock()
	err = f.readLine(fn)
	f.mux.Unlock()
	return
}

// ReadLines will return all lines in the form of an a bytes.Buffer
// Provided function can return true to end iteration early
func (f *File) ReadLines(fn func(*bytes.Buffer) bool) (err error) {
	var end bool
	f.mux.Lock()
	for err == nil && !end {
		err = f.readLine(func(b *bytes.Buffer) {
			end = fn(b)
		})
	}
	f.mux.Unlock()

	if err == io.EOF {
		err = nil
	}
	return
}

// Location will return the location to this file
func (f *File) Location() (loc string) {
	f.mux.Lock()
	loc = f.fLoc
	f.mux.Unlock()
	return
}

// Close will close the File
func (f *File) Close() (err error) {
	f.mux.Lock()
	if f.closed {
		err = ErrIsClosed
		goto END
	}

	if err = f.buf.Flush(); err != nil {
		return
	}

	if err = f.f.Close(); err != nil {
		return
	}

	f.f = nil
	f.buf = nil
	f.closed = true

END:
	f.mux.Unlock()
	return
}