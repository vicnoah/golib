package appdata

import (
	"errors"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"git.sabertrain.com/vector-dev/golib/helper/vlog"
	"git.sabertrain.com/vector-dev/golib/os/path/vhome"
)

const (
	_LOG_PATH = "log"
)

var once sync.Once

// SETUP 初始化日志
func SETUP(appName string) {
	once.Do(func() {
		vlog.SETUP(&LogFileWriter{
			logPath: path.Join(vhome.UserAppDataPath(), appName, _LOG_PATH),
		})
	})
}

// LogFileWriter 本地日志接口
type LogFileWriter struct {
	sync.Mutex
	file    *os.File
	size    int64
	logPath string
}

// Write 写入日志
func (p *LogFileWriter) Write(data []byte) (n int, err error) {
	p.Lock()
	defer p.Unlock()
	if p == nil {
		err = errors.New("logFileWriter is nil")
		return
	}
	if p.file == nil {
		er := p.createLog()
		if er != nil {
			err = er
			return
		}
	}
	n, err = p.file.Write(data)
	p.size += int64(n)
	//文件最大 64 K byte
	if p.size > 1024*64 {
		p.file.Close()
		er := p.createLog()
		if er != nil {
			err = er
			return
		}
	}
	return
}

// Close 关闭日志处理对象
func (p *LogFileWriter) Close() error {
	return p.file.Close()
}

func (p *LogFileWriter) createLog() (err error) {
	p.file, err = os.OpenFile(path.Join(p.logPath, strconv.FormatInt(time.Now().Unix(), 10)+".log"), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	p.size = 0
	return
}
