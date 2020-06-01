package toml

import (
	"github.com/BurntSushi/toml"
)

// Decode decode toml to struct
func Decode(data string, v interface{}) (err error) {
	_, err = toml.Decode(data, v)
	return
}
