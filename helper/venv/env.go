package venv

import (
	"os"
	"strings"
)

// Get 获取环境变量
func Get(key string) (value string) {
	for _, ev := range os.Environ() {
		kv := strings.Split(ev, "=")
		k := kv[0]
		v := kv[1]
		if k == key {
			value = v
			return
		}
	}
	return
}
