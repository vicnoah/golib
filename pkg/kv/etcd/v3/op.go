package v3

import (
	"go.etcd.io/etcd/clientv3"
)

// OpOption 操作配置
type OpOption clientv3.OpOption

func optsToV3(opts []OpOption) (v3opts []clientv3.OpOption) {
	for _, v := range opts {
		v3opts = append(v3opts, clientv3.OpOption(v))
	}
	return
}

// Op 操作
type Op clientv3.Op

// LeaseID lease id
type LeaseID clientv3.LeaseID

// WithLease attaches a lease ID to a key in 'Put' request.
func WithLease(leaseID LeaseID) OpOption {
	return OpOption(clientv3.WithLease(clientv3.LeaseID(leaseID)))
}

// WithPrefix watch prefix
func WithPrefix() OpOption {
	return OpOption(clientv3.WithPrefix())
}
