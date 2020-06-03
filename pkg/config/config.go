package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"git.sabertrain.com/vector-dev/golib/helper/venv"

	"git.sabertrain.com/vector-dev/golib/pkg/encoding/toml"
	"github.com/fsnotify/fsnotify"
)

const (
	// ENV_EXEC 执行所处环境,数据为空默认按照k8s处理,有数据直接按linux处理
	ENV_EXEC = "EXEC_ENV"
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
			configFiles: make(map[string]string),
			kv:          make(map[string]string),
		}
	})
	return c
}

// NewToStruct 创建toStruct
func (c *CM) NewToStruct(key string, conf interface{}) func() {
	return func() {
		c.toStruct(key, conf)
	}
}

// ToStruct 类型
type ToStruct func(key string, conf interface{}) error

// CM 配置管理
type CM struct {
	configFiles map[string]string // 监听配置文件
	watches     []watch           // 配置热更新
	kv          map[string]string // fsnotify监控的真实文件地址与原文件的映射（为兼容软链接）
	mu          sync.RWMutex      // 锁
}

// watch 配置热重载
type watch struct {
	configFileKey string // 配置文件key
	onWatch       func() // 对应watch回调
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
	err = c.toStruct(key, conf)
	return
}

// AddWatch 添加观察文件
func (c *CM) AddWatch(key string, f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.watches = append(c.watches, watch{key, f})
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
				fmt.Println(event.Op.String())
				if !ok {
					return
				}
				if c.isk8s() {
					// kubernetes中
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						watcher.Remove(event.Name)
						watcher.Add(c.kv[event.Name])
						c.notifyChange(event.Name)
						delete(c.kv, event.Name)
					}
				} else {
					// linux中
					if event.Op&fsnotify.Write == fsnotify.Write {
						c.notifyChange(event.Name)
						continue
					}
					// 为兼容vim修改,不完美vim使用:wq保存正常，使用:w会造成文件监控失败同时不能够收到后续通知
					if event.Op&fsnotify.Rename == fsnotify.Rename {
						// 文件修改事件中实际会产生Rename事件对原文件监控会失效，需要先移除之前的监控，然后添加新的监控
						watcher.Remove(event.Name)
						watcher.Add(event.Name)
						c.notifyChange(event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				fmt.Println(err)
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = c.loadWatcher(watcher)
	if err != nil {
		log.Fatal(err)
	}
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
	fmt.Println("配置更新", watch.configFileKey)
	watch.onWatch()
}

// findWatch 寻找监听器
func (c *CM) findWatch(fileName string) (watch, bool) {
	fileName = c.kv[fileName]
	for _, v := range c.watches {
		filePath := c.configFiles[v.configFileKey]
		if fileName == filePath {
			return v, true
		}
	}
	return watch{}, false
}

// toStruct 转换数据
func (c *CM) toStruct(key string, conf interface{}) (err error) {
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

// loadWatcher 重载监听器
func (c *CM) loadWatcher(watcher *fsnotify.Watcher) (err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i, v := range c.watches {
		if i == 0 {
			filePath, ok := c.configFiles[v.configFileKey]
			if !ok {
				err = ErrConfigNotFound
				return
			}
			fileName, _ := filepath.EvalSymlinks(filePath)
			// 存储路径
			c.kv[fileName] = filePath
			if err != nil {
				log.Fatal(err)
			}
			err = watcher.Add(fileName)
			if err != nil {
				return
			}
		}
	}
	return
}

// isk8s 是否在k8s中运行
func (c *CM) isk8s() bool {
	return venv.Get(ENV_EXEC) == ""
}
