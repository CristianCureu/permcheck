package internal

import (
	"os"
	"sync"
)

func FixPermissions(path string, mode os.FileMode, fixMode bool) error {
	if !fixMode {
		return nil
	}

	return os.Chmod(path, mode)
}

func WithLock(lock *sync.Mutex, fn func()) {
	lock.Lock()
	defer lock.Unlock()
	fn()
}
