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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ProcessQpsTokenReq(t *testing.T) {
	assert := assert.New(t)
	qps := (adapters[RATELIMITER])()
	b := qps.TryAccept("function", 100, 10)
	assert.True(b)
	b = qps.TryAccept("function.Schema1.op1", 10, 10)
	assert.True(b)
}

func Test_UpdateRateLimit(t *testing.T) {
	l := (adapters[RATELIMITER])()
	l.UpdateRateLimit("function.api1.limit", 200, 1)
	l.UpdateRateLimit("function.api1.limit", 100, 1)
}

func Test_DeleteRateLimit(t *testing.T) {
	qps := (adapters[RATELIMITER])()
	qps.DeleteRateLimiter("function.api1.limit")
}

func Test_RateLimitersTryAccept(t *testing.T) {
	assert := assert.New(t)
	after := time.After(1 * time.Second)
	count := 0
	stop := false
	for !stop {
		select {
		case <-after:
			fmt.Println(count)             //接近100
			assert.InDelta(count, 80, 120) //80~120
			stop = true
		default:
			pass := (adapters[RATELIMITER])().TryAccept("dongjiang", 100, 2)
			if pass {
				count++
			}
		}
	}
}
