package vcommon

import (
	"context"
)

// HandleTimeout 通用超时处理
func HandleTimeout(ctx context.Context, f func(), err *error) {
	// 开始工作负载前检查超时
	select {
	case <-ctx.Done():
		*err = ctx.Err()
		return
	default:
	}

	done := make(chan int)

	// 创建协程执行函数
	go func() {
		defer func() {
			close(done)
		}()
		f()
		done <- -1
	}()

	select {
	case <-ctx.Done():
		*err = ctx.Err()
		return
	case <-done:
		return
	}
}
