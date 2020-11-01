package v3

import "time"

var testConfig = Config{
	Endpoints:   []string{"127.0.0.1:2379"},
	DialTimeout: time.Second * 5,
}
