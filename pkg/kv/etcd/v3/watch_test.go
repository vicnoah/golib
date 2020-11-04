package v3

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	cli, err := New(testConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to etcd success.")
	defer cli.Close()

	go func() {
		// watch key:q1mi change
		rch := cli.Watch(context.Background(), "q1mi") // <-chan WatchResponse
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()
	gr, err := cli.Grant(context.TODO(), 1)
	if err != nil {
		log.Fatal(err)
	}
	_, err = cli.Put(context.TODO(), "q1mi", "dsb", WithLease(LeaseID(gr.ID)))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	_, err = cli.Delete(context.TODO(), "q1mi")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i <= 5; i++ {
		_, err = cli.Put(context.TODO(), "q1mi", "dsb")
		if err != nil {
			log.Fatal(err)
		}
	}
}
