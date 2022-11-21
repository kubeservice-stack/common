package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	assert := assert.New(t)

	var cases = []struct {
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

	var cases = []struct {
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
