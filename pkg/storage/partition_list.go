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

package storage

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// partitionList represents a linked list for partitions.
// Each partition is arranged in order order of newest to oldest.
type partitionList interface {
	insert(partition partition)
	remove(partition partition) error
	swap(old, new partition) error
	getHead() partition
	size() int
	newIterator() partitionIterator
	String() string
}

//
// for iterator.next() {
//   partition, err := iterator.value()
//      //Do something with partition
// }
// Iterator represents an iterator for partition list.
type partitionIterator interface {
	next() bool
	value() partition
	currentNode() *partitionNode
}

type partitionListImpl struct {
	numPartitions int64
	head          *partitionNode
	tail          *partitionNode
	mu            sync.RWMutex
}

func newPartitionList() partitionList {
	return &partitionListImpl{}
}

func (p *partitionListImpl) getHead() partition {
	if p.size() <= 0 {
		return nil
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.head.value()
}

func (p *partitionListImpl) insert(partition partition) {
	node := &partitionNode{
		val: partition,
	}
	p.mu.RLock()
	head := p.head
	p.mu.RUnlock()
	if head != nil {
		node.next = head
	}

	p.setHead(node)
	atomic.AddInt64(&p.numPartitions, 1)
}

func (p *partitionListImpl) remove(target partition) error {
	if p.size() <= 0 {
		return fmt.Errorf("empty partition")
	}

	// Iterate over itself from the head.
	var prev, next *partitionNode
	iterator := p.newIterator()
	for iterator.next() {
		current := iterator.currentNode()
		if !samePartitions(current.value(), target) {
			prev = current
			continue
		}

		// remove the current node.

		iterator.next()
		next = iterator.currentNode()
		switch {
		case prev == nil:
			// removing the head node
			p.setHead(next)
		case next == nil:
			// removing the tail node
			prev.setNext(nil)
			p.setTail(prev)
		default:
			// removing the middle node
			prev.setNext(next)
		}
		atomic.AddInt64(&p.numPartitions, -1)

		if err := current.value().clean(); err != nil {
			return fmt.Errorf("failed to clean resources managed by partition to be removed: %w", err)
		}
		return nil
	}

	return fmt.Errorf("the given partition was not found")
}

func (p *partitionListImpl) swap(old, new partition) error {
	if p.size() <= 0 {
		return fmt.Errorf("empty partition")
	}

	// Iterate over itself from the head.
	var prev, next *partitionNode
	iterator := p.newIterator()
	for iterator.next() {
		current := iterator.currentNode()
		if !samePartitions(current.value(), old) {
			prev = current
			continue
		}

		// swap the current node.

		newNode := &partitionNode{
			val:  new,
			next: current.getNext(),
		}
		iterator.next()
		next = iterator.currentNode()
		switch {
		case prev == nil:
			// swapping the head node
			p.setHead(newNode)
		case next == nil:
			// swapping the tail node
			prev.setNext(newNode)
			p.setTail(newNode)
		default:
			// swapping the middle node
			prev.setNext(newNode)
		}
		return nil
	}

	return fmt.Errorf("the given partition was not found")
}

func samePartitions(x, y partition) bool {
	return x.minTimestamp() == y.minTimestamp()
}

func (p *partitionListImpl) size() int {
	return int(atomic.LoadInt64(&p.numPartitions))
}

func (p *partitionListImpl) newIterator() partitionIterator {
	p.mu.RLock()
	head := p.head
	p.mu.RUnlock()
	// Put a dummy node so that it positions the head on the first next() call.
	dummy := &partitionNode{
		next: head,
	}
	return &partitionIteratorImpl{
		current: dummy,
	}
}

func (p *partitionListImpl) setHead(node *partitionNode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.head = node
}

func (p *partitionListImpl) setTail(node *partitionNode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.tail = node
}

func (p *partitionListImpl) String() string {
	b := &strings.Builder{}
	iterator := p.newIterator()
	for iterator.next() {
		p := iterator.value()
		if _, ok := p.(*memoryPartition); ok {
			b.WriteString("[Memory Partition]")
		} else {
			b.WriteString("[Unknown Partition]")
		}
		b.WriteString("->")
	}
	return strings.TrimSuffix(b.String(), "->")
}

// partitionNode wraps a partition to hold the pointer to the next one.
type partitionNode struct {
	// val is immutable
	val  partition
	next *partitionNode
	mu   sync.RWMutex
}

// value gives back the actual partition of the node.
func (p *partitionNode) value() partition {
	return p.val
}

func (p *partitionNode) setNext(node *partitionNode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.next = node
}

func (p *partitionNode) getNext() *partitionNode {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.next
}

type partitionIteratorImpl struct {
	current *partitionNode
}

func (i *partitionIteratorImpl) next() bool {
	if i.current == nil {
		return false
	}
	next := i.current.getNext()
	i.current = next
	return i.current != nil
}

func (i *partitionIteratorImpl) value() partition {
	if i.current == nil {
		return nil
	}
	return i.current.value()
}

func (i *partitionIteratorImpl) currentNode() *partitionNode {
	return i.current
}
