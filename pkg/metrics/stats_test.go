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

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		hit  int
		all  int
		rate float64
	}{
		{3, 4, 0.75},
		{0, 1, 0.0},
		{3, 3, 1.0},
		{0, 0, 0.0},
	}

	for _, cs := range cases {
		st := &Stats{}
		for i := 0; i < cs.hit; i++ {
			st.IncrHitCount()
		}
		for i := 0; i < cs.all; i++ {
			st.IncrAllCount()
		}
		assert.Equal(cs.rate, st.HitRate(), "not equal")
	}
}

func TestStatsAllHitCount(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		hit  int
		all  int
		rate float64
	}{
		{3, 4, 0.75},
	}
	for _, cs := range cases {
		st := &Stats{}
		for i := 0; i < cs.hit; i++ {
			st.IncrHitCount()
		}
		for i := 0; i < cs.all; i++ {
			st.IncrAllCount()
		}
		assert.Equal(cs.rate, st.HitRate(), "not equal")
		assert.Equal(st.AllCount(), uint64(cs.all), "is not equal")
	}
}
