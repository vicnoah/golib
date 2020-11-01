package v3

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestLock(t *testing.T) {
	cli, err := New(testConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to etcd success.")
	defer cli.Close()

	m1, s1, err := cli.NewMutex("/my-lock/")
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()

	m2, s2, err := cli.NewMutex("/my-lock/")
	if err != nil {
		log.Fatal(err)
	}
	defer s2.Close()

	// 会话s1获取锁
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("acquired lock for s1")

	m2Locked := make(chan struct{})
	go func() {
		defer close(m2Locked)
		// 等待直到会话s1释放了/my-lock/的锁
		if err := m2.Lock(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("released lock for s1")

	<-m2Locked
	fmt.Println("acquired lock for s2")
}
