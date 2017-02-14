package lineFile

import (
	"os"

	"github.com/itsmontoya/async/file"
)

// getOS is an OpenFile func for the os backend
func getOS(name string, flag int, perm os.FileMode) (fi FileInt, err error) {
	return os.OpenFile(name, flag, perm)
}

// getAIO is an OpenFile func for the aio backend
func getAIO(name string, flag int, perm os.FileMode) (fi FileInt, err error) {
	return file.OpenFile(name, flag, perm)
}
