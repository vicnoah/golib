package vpanic

import (
	"fmt"
	"runtime"
	"time"
)

// PanicWrapper Prevent possible panic errors
func PanicWrapper(f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("internal error: %v\n", err)
			func() {
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				fmt.Println(string(buf[:n]))
			}()
			time.Sleep(time.Second * 10)
			PanicWrapper(f)
		}
	}()
	f()
}
