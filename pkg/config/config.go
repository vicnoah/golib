package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"git.sabertrain.com/vector-dev/golib/pkg/encoding/toml"
	"github.com/fsnotify/fsnotify"
)

var (
	// ErrConfigNotFound 配置未找到错误
	ErrConfigNotFound = errors.New("config file not found")
)

var (
	c    *CM
	once sync.Once
)

// New 新建配置管理器
func New() *CM {
	once.Do(func() {
		c = &CM{
			configFiles: make(map[string]string, 0),
		}
	})
	return c
}

// CM 配置管理
type CM struct {
	configFiles map[string]string // 监听配置文件
	watches     []watch           // 配置热更新
	mu          sync.RWMutex      // 锁
}

// watch 配置热重载
type watch struct {
	configFileKey string                 // 配置文件key
	conf          interface{}            // 配置结构体地址
	onWatch       func(conf interface{}) // 对应watch回调
}

// Add 添加配置文件
func (c *CM) Add(key, filePath string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.configFiles[key] = filePath
}

// ToStruct 转换配置到结构体
func (c *CM) ToStruct(key string, conf interface{}) (err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	filePath, ok := c.configFiles[key]
	if !ok {
		err = ErrConfigNotFound
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	dataBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = toml.Decode(string(dataBytes), conf)
	return
}

// AddWatch 添加观察文件
func (c *CM) AddWatch(key string, conf interface{}, f func(conf interface{})) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cpConf := conf
	c.watches = append(c.watches, watch{key, cpConf, f})
}

// StartWatch 开始热更新
func (c *CM) StartWatch() {
	go c.watch()
}

// watch 监听配置
func (c *CM) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					c.notifyChange(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	c.mu.RLock()
	for _, v := range c.watches {
		filePath, ok := c.configFiles[v.configFileKey]
		if !ok {
			err = ErrConfigNotFound
			return
		}
		err = watcher.Add(filePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	c.mu.RUnlock()
	<-done
}

// notifyChange 通知改变
func (c *CM) notifyChange(fileName string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	watch, ok := c.findWatch(fileName)
	if !ok {
		return
	}
	c.ToStruct(watch.configFileKey, watch.conf)
	watch.onWatch(watch.conf)
}

// findWatch 寻找监听器
func (c *CM) findWatch(fileName string) (watch, bool) {
	for _, v := range c.watches {
		filePath := c.configFiles[v.configFileKey]
		if fileName == filePath {
			return v, true
		}
	}
	return watch{}, false
}
