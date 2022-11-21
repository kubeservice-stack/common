package ratelimiter

import (
	"sync"

	"k8s.io/client-go/util/flowcontrol"
)

const (
	RATELIMITER = "RATELIMITER"
)

type RateLimiters struct {
	sync.RWMutex //协程 安全
	m            map[string]flowcontrol.RateLimiter
}

var (
	once       = new(sync.Once)
	qpsLimiter *RateLimiters
)

func NewRateLimiters() Limiter {
	once.Do(func() {
		qpsLimiter = &RateLimiters{m: make(map[string]flowcontrol.RateLimiter)}
	})
	return qpsLimiter
}

func (l *RateLimiters) TryAccept(name string, qps, burst int) bool {
	l.RLock()
	limiter, ok := l.m[name]
	if !ok {
		l.RUnlock()
		return l.addLimiter(name, qps, burst) //新增
	}
	l.RUnlock()
	return limiter.TryAccept()

}

func (l *RateLimiters) addLimiter(name string, qps, burst int) bool {
	var bucketSize int
	if qps >= 1 {
		bucketSize = qps
	} else {
		bucketSize = DefaultRate
	}
	l.Lock()
	// 新建token bucket
	r := flowcontrol.NewTokenBucketRateLimiter(float32(bucketSize), burst)
	l.m[name] = r
	l.Unlock()
	return r.TryAccept()
}

func (l *RateLimiters) UpdateRateLimit(name string, qps, burst int) {
	l.addLimiter(name, qps, burst)
}

func (l *RateLimiters) DeleteRateLimiter(name string) {
	l.Lock()
	delete(l.m, name)
	l.Unlock()
}

func init() {
	Register(RATELIMITER, NewRateLimiters)
}
