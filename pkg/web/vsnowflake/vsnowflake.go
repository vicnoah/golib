// Package vsnowflake 雪花算法生成数据唯一id
//内存安全的唯一id生成器
package vsnowflake

import (
	"errors"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	once sync.Once
	sn   *SnowFlake
)

func init() {
	once.Do(func() {
		sn = &SnowFlake{}
	})
}

// SnowFlake 雪花算法对象
type SnowFlake struct {
	mu sync.Mutex
}

// New 返回一个64位雪花ID
func (s *SnowFlake) New() (uint64, error) {
	defer s.mu.Unlock()
	s.mu.Lock()
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		//fmt.Println(err)
		return 0, errors.New("generate id error")
	}

	// Generate a snowflake ID.
	//id := node.Generate()

	// Print out the ID in a few different ways.
	//fmt.Printf("Int64  ID: %d\n", id)
	//fmt.Printf("String ID: %s\n", id)
	//fmt.Printf("Base2  ID: %s\n", id.Base2())
	//fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	//fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	//fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	//fmt.Printf("ID Step  : %d\n", id.Step())

	// Generate and print, all in one.
	//fmt.Printf("ID       : %d\n", node.Generate().Int64())
	return uint64(node.Generate().Int64()), nil
}

// Get 返回一个32位雪花ID
func Get() (uint64, error) {
	//延迟时间以保证id唯一
	time.Sleep(time.Millisecond * 1)
	return sn.New()
}
