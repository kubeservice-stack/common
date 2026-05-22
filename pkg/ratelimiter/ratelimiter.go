/*
Copyright 2022 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ratelimiter

import (
	"sync"

	"k8s.io/client-go/util/flowcontrol"
)

const (
	RATELIMITER = "RATELIMITER"
)

type RateLimiters struct {
	sync.Mutex         // protects m
	m      map[string]flowcontrol.RateLimiter
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
	limiter := l.getOrCreateLimiter(name, qps, burst)
	return limiter.TryAccept()
}

func (l *RateLimiters) getOrCreateLimiter(name string, qps, burst int) flowcontrol.RateLimiter {
	l.Lock()
	defer l.Unlock()
	if limiter, ok := l.m[name]; ok {
		return limiter
	}
	l.m[name] = newRateLimiter(qps, burst)
	return l.m[name]
}

func newRateLimiter(qps, burst int) flowcontrol.RateLimiter {
	var bucketSize int
	if qps >= 1 {
		bucketSize = qps
	} else {
		bucketSize = DefaultRate
	}
	return flowcontrol.NewTokenBucketRateLimiter(float32(bucketSize), burst)
}

// addLimiter is required by Limiter interface; use UpdateRateLimit for external callers
func (l *RateLimiters) addLimiter(name string, qps, burst int) bool {
	l.Lock()
	l.m[name] = newRateLimiter(qps, burst)
	l.Unlock()
	return true
}

func (l *RateLimiters) UpdateRateLimit(name string, qps, burst int) {
	l.Lock()
	l.m[name] = newRateLimiter(qps, burst)
	l.Unlock()
}

func (l *RateLimiters) DeleteRateLimiter(name string) {
	l.Lock()
	delete(l.m, name)
	l.Unlock()
}

func init() {
	Register(RATELIMITER, NewRateLimiters)
}
