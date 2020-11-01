package v3

import (
	"go.etcd.io/etcd/clientv3/concurrency"
)

// NewMutex 互斥锁
// 使用完锁需调用se.Close()关闭会话
func (ec *Etcd) NewMutex(prx string) (mu *concurrency.Mutex, se *concurrency.Session, err error) {
	se, err = concurrency.NewSession(ec.ocli)
	if err != nil {
		return
	}
	mu = concurrency.NewMutex(se, prx)
	return
}
