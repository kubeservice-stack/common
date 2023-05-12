/*
Copyright 2023 The KubeService-Stack Authors.

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

package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Timer(t *testing.T) {
	assert := assert.New(t)
	tp := NewTimerPool()
	a := tp.Get(time.Second)
	tp.Put(a)
	b := tp.Get(time.Second)
	assert.NotNil(a)
	assert.Equal(a, b)
}

func Test_TimerTimeout(t *testing.T) {
	assert := assert.New(t)
	tp := NewTimerPool()
	a := time.NewTimer(10 * time.Second)
	tp.Put(a)
	tp.Put(a)
	tp.Put(a)
	assert.NotEmpty(tp.Get(10 * time.Second))
}
