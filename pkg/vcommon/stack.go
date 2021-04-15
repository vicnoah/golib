package vcommon

import "runtime"

// Stack 返回最近错误的调用堆栈
func Stack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return "==> " + string(buf[:n]) + "\n"
}
