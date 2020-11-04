package v3

import "time"

var testConfig = Config{
	Endpoints:   []string{"192.168.0.25:2379"},
	DialTimeout: time.Second * 5,
}
