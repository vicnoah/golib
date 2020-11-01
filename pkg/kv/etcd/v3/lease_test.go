package v3

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestLease(t *testing.T) {
	cli, err := New(testConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to etcd success.")
	defer cli.Close()

	// 创建一个5秒的租约
	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
	}

	// 5秒钟之后, /nazha/ 这个key就会被移除
	_, err = cli.Put(context.TODO(), "/nazha/", "dsb", WithLease(LeaseID(resp.ID)))
	if err != nil {
		log.Fatal(err)
	}
}
