package metrics

import (
	"sync/atomic"
)

type Stats struct {
	hitCount uint64 // hit count
	allCount uint64 // not hit count = all count - hit count
}

func (st *Stats) IncrHitCount() uint64 {
	return atomic.AddUint64(&st.hitCount, 1)
}

func (st *Stats) IncrAllCount() uint64 {
	return atomic.AddUint64(&st.allCount, 1)
}

func (st *Stats) HitCount() uint64 {
	return atomic.LoadUint64(&st.hitCount)
}

func (st *Stats) AllCount() uint64 {
	return atomic.LoadUint64(&st.allCount)
}

// HitRate returns rate for cache hitting
func (st *Stats) HitRate() float64 {
	hc := st.HitCount()
	total := st.AllCount()
	if total == 0 {
		return 0.0
	}
	return float64(hc) / float64(total)
}
