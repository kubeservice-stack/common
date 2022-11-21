package discovery

import (
	"fmt"

	"github.com/pkg/errors"
	etcdcliv3 "go.etcd.io/etcd/clientv3"
)

var (
	ErrWatchFailed = fmt.Errorf("etcd watch returns a nil chan")
	ErrNoKey       = fmt.Errorf("etcd has no such key")
	ErrTxnFailed   = fmt.Errorf("role changed or target revision mismatch")
	ErrTxnConvert  = fmt.Errorf("cannot covert etcd transaction")
)

// 将txn响应和错误转化为一个错误
func TxnErr(resp *etcdcliv3.TxnResponse, err error) error {
	if err != nil {
		return errors.WithStack(err)
	}
	if !resp.Succeeded {
		return ErrTxnFailed
	}
	return nil
}
