package vlog

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	once   sync.Once
	logger *logrus.Logger
	wc     io.WriteCloser
)

// SETUP 初始化日志库
func SETUP(iwc io.WriteCloser) {
	once.Do(func() {
		logger = logrus.New()
		wc = iwc
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetOutput(iwc)
		logger.SetLevel(logrus.ErrorLevel)
	})
}

// Do 获取日志对象
func Do() *logrus.Logger {
	return logger
}

// Close 关闭对象
func Close() error {
	return wc.Close()
}
