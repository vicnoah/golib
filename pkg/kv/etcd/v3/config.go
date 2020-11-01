package v3

import "time"

// Config 配置管理
type Config struct {
	// Endpoints etcd端点
	Endpoints []string
	// DialTimeout 超时
	DialTimeout time.Duration
}
