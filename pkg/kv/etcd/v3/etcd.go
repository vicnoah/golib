package v3

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// New 新建etcd连接
func New(cfg Config) (ec *Etcd, err error) {
	cli, err := clientv3.New(clientv3.Config(cfg))
	if err != nil {
		return
	}
	ec = &Etcd{
		cli: cli,
	}
	return
}

// Etcd etcd实例
type Etcd struct {
	cli *clientv3.Client
}

// Close 关闭客户端
func (ec *Etcd) Close() error {
	return ec.cli.Close()
}

// PutResponse put响应
type PutResponse clientv3.PutResponse

// Put 提交数据
func (ec *Etcd) Put(ctx context.Context, key, val string, opts ...OpOption) (rs *PutResponse, err error) {
	ors, err := ec.cli.Put(ctx, key, val, optsToV3(opts)...)
	if err != nil {
		return
	}
	rsObj := PutResponse(*ors)
	rs = &rsObj
	return
}

// GetResponse get响应
type GetResponse clientv3.GetResponse

// Get 查询数据
func (ec *Etcd) Get(ctx context.Context, key string, opts ...OpOption) (rs *GetResponse, err error) {
	ors, err := ec.cli.Get(ctx, key, optsToV3(opts)...)
	if err != nil {
		return
	}
	rsObj := GetResponse(*ors)
	rs = &rsObj
	return
}

// DeleteResponse 删除响应
type DeleteResponse clientv3.DeleteResponse

// Delete 删除数据
func (ec *Etcd) Delete(ctx context.Context, key string, opts ...OpOption) (rs *DeleteResponse, err error) {
	ors, err := ec.cli.Delete(ctx, key, optsToV3(opts)...)
	if err != nil {
		return
	}
	rsObj := DeleteResponse(*ors)
	rs = &rsObj
	return
}

// LeaseGrantResponse lease grant response
type LeaseGrantResponse clientv3.LeaseGrantResponse

// Grant 创建过期租约
// 单位秒
func (ec *Etcd) Grant(ctx context.Context, ttl int64) (rs *LeaseGrantResponse, err error) {
	ors, err := ec.cli.Grant(ctx, ttl)
	if err != nil {
		return
	}
	rsObj := LeaseGrantResponse(*ors)
	rs = &rsObj
	return
}

// LeaseKeepAliveResponse lease保持活跃
type LeaseKeepAliveResponse clientv3.LeaseKeepAliveResponse

// KeepAlive 保持活跃
func (ec *Etcd) KeepAlive(ctx context.Context, id LeaseID) (rs <-chan *LeaseKeepAliveResponse, err error) {
	rsw := make(chan *LeaseKeepAliveResponse)
	ors, err := ec.cli.KeepAlive(ctx, clientv3.LeaseID(id))
	if err != nil {
		return
	}
	go func() {
		defer func() {
			close(rsw)
		}()
		for ka := range ors {
			kar := LeaseKeepAliveResponse(*ka)
			rsw <- &kar
		}
	}()
	rs = rsw
	return
}

// WatchChan 侦听channel
type WatchChan clientv3.WatchChan

// Watch 监听
func (ec *Etcd) Watch(ctx context.Context, key string, opts ...OpOption) (watchChan WatchChan) {
	ors := ec.cli.Watch(ctx, key, optsToV3(opts)...)
	watchChan = WatchChan(ors)
	return
}
