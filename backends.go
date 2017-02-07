package lineFile

import (
	"os"

	"github.com/itsmontoya/aio"
)

// TODO: Figure out a better way to handle initializing aio
var am *aio.AIO

// getOS is an OpenFile func for the os backend
func getOS(name string, flag int, perm os.FileMode) (fi FileInt, err error) {
	return os.OpenFile(name, flag, perm)
}

// getAIO is an OpenFile func for the aio backend
func getAIO(name string, flag int, perm os.FileMode) (fi FileInt, err error) {
	if am == nil {
		am = aio.New(1)
	}

	return am.OpenFile(name, flag, perm)
}
