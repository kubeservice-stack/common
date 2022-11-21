package item

import (
	"container/list"
	"time"
)

type LfuItem struct {
	Key         interface{}
	Value       interface{}
	FreqElement *list.Element
	Expiration  *time.Time
}

// returns boolean value whether this item is expired or not.
func (it *LfuItem) IsExpired(now *time.Time) bool {
	if it.Expiration == nil {
		return false
	}
	if now == nil {
		t := time.Now()
		now = &t
	}
	return it.Expiration.Before(*now)
}

func (it *LfuItem) Expire() *time.Time {
	return it.Expiration
}
