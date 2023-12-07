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

package tokenbucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenBucket(t *testing.T) {
	type testCase struct {
		Name       string
		Rate       uint64
		TimeWindow time.Duration
		Equal      bool
	}

	testcases := []testCase{
		{
			Name:       "New(0, 0)",
			Rate:       0,
			TimeWindow: 0,
		},
		{
			Name:       "New(0, 1)",
			Rate:       0,
			TimeWindow: 1,
		},
		{
			Name:       "New(1, 1)",
			Rate:       1,
			TimeWindow: 1,
		},
		{
			Name:       "New(1, 10)",
			Rate:       1,
			TimeWindow: 10,
		},
		{
			Name:       "New(1, 0)",
			Rate:       1,
			TimeWindow: 0,
		},
	}

	for _, tc := range testcases {
		tb := New(tc.Name, tc.Rate, tc.TimeWindow*time.Second)
		for i := 0; i < 10; i++ {
			if i == 0 {
				require.Equal(t, false, tb.Limit(), tc.Name)
			} else {
				require.Equal(t, true, tb.Limit(), tc.Name)
			}
		}
	}

}

func TestUnixNano(t *testing.T) {
	assert := assert.New(t)
	a := unixNano()
	assert.GreaterOrEqual(uint64(time.Now().UnixNano()), a)
}

func TestUpdateCapacity(t *testing.T) {
	assert := assert.New(t)
	tb := New("TestUpdateRate", 4, 2*time.Second)
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), false)
	tb.UpdateCapacity(3)
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), true)
	assert.Equal(tb.Limit(), true)
	time.Sleep(2 * time.Second)
	// 更新rate下一个时间窗口生效
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), true)
	assert.Equal(tb.Limit(), true)
}

func TestUndo(t *testing.T) {
	assert := assert.New(t)
	tb := New("TestUndo", 3, 1*time.Second)
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), true)
	tb.Undo()
	assert.Equal(tb.Limit(), false)
	assert.Equal(tb.Limit(), true)
	assert.Equal(tb.Limit(), true)
}
