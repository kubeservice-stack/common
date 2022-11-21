package ratelimiter

import (
	"fmt"
)

var (
	ErrLimiterRegisterAdapterNil    = fmt.Errorf("Limiter: Register adapter is nil")
	ErrLimiterDoubleRegisterAdapter = fmt.Errorf("Limiter: Register called twice for adapter")
)

const (
	DefaultRate = 2147483647
)

type Limiter interface {
	TryAccept(name string, qps, burst int) bool  // Accept请求， limiter接受正常返回，不接受返回false
	addLimiter(name string, qps, burst int) bool // 添加一组limiter, 添加成功true, name 重复 false
	UpdateRateLimit(name string, qps, burst int) // 更新qps, 并发burst接口，etcd配置变化
	DeleteRateLimiter(name string)               // 清理 limiter
}

type Instance func() Limiter

var adapters = make(map[string]Instance)

func Register(name string, adapter Instance) {
	if adapter == nil {
		panic(ErrLimiterRegisterAdapterNil.Error())
	}
	if _, ok := adapters[name]; ok {
		panic(ErrLimiterDoubleRegisterAdapter.Error() + name)
	}
	adapters[name] = adapter
}

func GetLimiter(name string) (Instance, bool) {
	if limiter, ok := adapters[name]; ok {
		return limiter, ok
	}
	return nil, false
}

func GetDefaultLimiter() Instance {
	return adapters[RATELIMITER]
}

func HasRegister(name string) bool {

	if _, ok := adapters[name]; ok {
		return true
	}
	return false
}
