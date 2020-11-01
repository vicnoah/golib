package v3

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestKeepAlive(t *testing.T) {
	cli, err := New(testConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connect to etcd success.")
	defer cli.Close()

	resp, err := cli.Grant(context.TODO(), 2)
	if err != nil {
		log.Fatal(err)
	}

	_, err = cli.Put(context.TODO(), "/nazha/", "dsb", WithLease(LeaseID(resp.ID)))
	if err != nil {
		log.Fatal(err)
	}

	// the key 'foo' will be kept forever
	ch, kaerr := cli.KeepAlive(context.TODO(), LeaseID(resp.ID))
	if kaerr != nil {
		log.Fatal(kaerr)
	}
	count := 0
	for {
		if count == 10 {
			break
		}
		ka := <-ch
		fmt.Println("ttl:", ka.TTL)
		count++
	}
}
