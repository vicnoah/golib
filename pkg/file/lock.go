package file

import (
	"fmt"
	"os"
	"syscall"
)

// Lock 文件锁
type Lock struct {
	dir string
}

// New 创建FileLock对象
func New(dir string) *Lock {
	return &Lock{
		dir: dir,
	}
}

// Lock 加锁
func (l *Lock) Lock(f *os.File) error {
	// |syscall.LOCK_NB
	err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		return fmt.Errorf("cannot flock %s - %s", l.dir, err)
	}
	return nil
}

// Unlock 释放锁
func (l *Lock) Unlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
