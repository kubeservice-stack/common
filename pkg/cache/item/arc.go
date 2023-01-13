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

package item

import (
	"container/list"
	"time"
)

type ArcItem struct {
	Key        interface{}
	Value      interface{}
	Expiration *time.Time
}

// returns boolean value whether this item is expired or not.
func (it *ArcItem) IsExpired(now *time.Time) bool {
	if it.Expiration == nil {
		return false
	}
	if now == nil {
		t := time.Now()
		now = &t
	}
	return it.Expire().Before(*now)
}

func (it *ArcItem) Expire() *time.Time {
	return it.Expiration
}

type ArcList struct {
	l    *list.List
	keys map[interface{}]*list.Element
}

func NewARCList() *ArcList {
	return &ArcList{
		l:    list.New(),
		keys: make(map[interface{}]*list.Element),
	}
}

// has key func
// return bool
func (al *ArcList) Has(key interface{}) bool {
	_, ok := al.keys[key]
	return ok
}

// Lookup func : search list.element for key
func (al *ArcList) Lookup(key interface{}) *list.Element {
	elt := al.keys[key]
	return elt
}

// Move item to front
func (al *ArcList) MoveToFront(elt *list.Element) {
	al.l.MoveToFront(elt)
}

// push item to front
func (al *ArcList) PushFront(key interface{}) {
	elt := al.l.PushFront(key)
	al.keys[key] = elt
}

// delete item
func (al *ArcList) Remove(key interface{}, elt *list.Element) {
	if al.Has(key) {
		delete(al.keys, key)
	}
	al.l.Remove(elt)
}

// delete last
func (al *ArcList) RemoveTail() interface{} {
	elt := al.l.Back()
	al.l.Remove(elt)

	key := elt.Value
	if al.Has(key) {
		delete(al.keys, key)
	}

	return key
}

// list len
func (al *ArcList) Len() int {
	return al.l.Len()
}
