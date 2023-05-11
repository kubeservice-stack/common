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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kubeservice-stack/common/pkg/logger"
	"github.com/kubeservice-stack/common/pkg/utils"
)

/* memory data tree. just for data code list
  │                 │
Read              Write
  │                 │
  │                 V
  │      ┌───────────────────┐ max: 1615010800
  ├─────>   Memory Partition
  │      └───────────────────┘ min: 1615007201
  │
  │      ┌───────────────────┐ max: 1615007200
  ├─────>   Memory Partition
  │      └───────────────────┘ min: 1615003601
  │
  │      ┌───────────────────┐ max: 1615003600
  └─────>   Memory Partition
         └───────────────────┘ min: 1615000000
*/
type Storage struct {
	partitionList partitionList

	partitionDuration  time.Duration
	retention          time.Duration
	timestampPrecision TimestampPrecision
	writeTimeout       time.Duration

	logger         logger.Logger
	workersLimitCh chan struct{}
	// be incremented to guarantee all writes are done gracefully.
	wg sync.WaitGroup
	// timerpool
	timerpool *utils.TimerPool

	doneCh chan struct{}
}

func NewStorage(opts ...Option) (StorageInterface, error) {
	s := &Storage{
		partitionList:      newPartitionList(),
		workersLimitCh:     make(chan struct{}, defaultWorkersLimit),
		partitionDuration:  defaultPartitionDuration,
		retention:          defaultRetention,
		timestampPrecision: defaultTimestampPrecision,
		writeTimeout:       defaultWriteTimeout,
		doneCh:             make(chan struct{}),
		timerpool:          utils.NewTimerPool(),
	}

	// setting option
	for _, opt := range opts {
		opt(s)
	}

	// new partition
	s.newPartition(nil)

	return s, nil
}

func (s *Storage) newPartition(p partition) error {
	if p == nil {
		p = NewMemoryPartition(s.partitionDuration, s.timestampPrecision)
	}
	s.partitionList.insert(p)
	return nil
}

func (s *Storage) InsertRows(rows []Row) error {
	s.wg.Add(1)
	defer s.wg.Done()

	insert := func() error {
		defer func() { <-s.workersLimitCh }()
		if err := s.ensureActiveHead(); err != nil {
			return err
		}
		iterator := s.partitionList.newIterator()
		n := s.partitionList.size()
		rowsToInsert := rows

		for i := 0; i < n && i < defaultwritablePartitionsNum; i++ {
			if len(rowsToInsert) == 0 {
				break
			}
			if !iterator.next() {
				break
			}
			outdatedRows, err := iterator.value().insertRows(rowsToInsert)
			if err != nil {
				return fmt.Errorf("failed to insert rows: %w", err)
			}
			rowsToInsert = outdatedRows
		}
		return nil
	}

	// Limit the number of concurrent goroutines to prevent from out of memory
	// errors and CPU trashing even if too many goroutines attempt to write.
	select {
	case s.workersLimitCh <- struct{}{}:
		return insert()
	default:
	}

	// Seems like all workers are busy; wait for up to writeTimeout

	t := s.timerpool.Get(s.writeTimeout)
	select {
	case s.workersLimitCh <- struct{}{}:
		s.timerpool.Put(t)
		return insert()
	case <-t.C:
		s.timerpool.Put(t)
		return fmt.Errorf("failed to write a data point in %s, since it is overloaded with %d concurrent writers",
			s.writeTimeout, defaultWorkersLimit)
	}
}

func (s *Storage) ensureActiveHead() error {
	head := s.partitionList.getHead()
	if head != nil && head.active() {
		return nil
	}

	// All partitions seems to be inactive so add a new partition to the list.
	if err := s.newPartition(nil); err != nil {
		return err
	}
	go func() {
		if err := s.flushPartitions(); err != nil {
			s.logger.Error("failed to flush in-memory partitions", logger.Error(err))
		}
	}()
	return nil
}

func (s *Storage) flushPartitions() error {
	i := 0
	iterator := s.partitionList.newIterator()
	for iterator.next() {
		if i < defaultwritablePartitionsNum {
			i++
			continue
		}
		part := iterator.value()
		if part == nil {
			return fmt.Errorf("unexpected empty partition found")
		}
		_, ok := part.(*memoryPartition)
		if !ok {
			continue
		}

		if err := s.partitionList.remove(part); err != nil {
			return fmt.Errorf("failed to remove partition: %w", err)
		}

	}
	return nil
}

func (s *Storage) Select(name string, labels []Label, start, end int64) ([]*DataPoint, error) {
	if name == "" {
		return nil, fmt.Errorf("metric must be set")
	}
	if start >= end {
		return nil, fmt.Errorf("the given start is greater than end")
	}
	points := make([]*DataPoint, 0)

	// Iterate over all partitions from the newest one.
	iterator := s.partitionList.newIterator()
	for iterator.next() {
		part := iterator.value()
		if part == nil {
			return nil, fmt.Errorf("unexpected empty partition found")
		}
		if part.minTimestamp() == 0 {
			// Skip the partition that has no points.
			continue
		}
		if part.maxTimestamp() < start {
			// No need to keep going anymore
			break
		}
		if part.minTimestamp() > end {
			continue
		}
		ps, err := part.selectDataPoints(name, labels, start, end)
		if errors.Is(err, ErrNoDataPoints) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to select data points: %w", err)
		}
		// in order to keep the order in ascending.
		points = append(ps, points...)
	}
	if len(points) == 0 {
		return nil, ErrNoDataPoints
	}
	return points, nil
}

func (s *Storage) Close() error {
	s.wg.Wait()
	close(s.doneCh)

	// TODO: Prevent from new goroutines calling InsertRows(), for graceful shutdown.

	// Make all writable partitions read-only by inserting as same number of those.
	for i := 0; i < defaultwritablePartitionsNum; i++ {
		if err := s.newPartition(nil); err != nil {
			return err
		}
	}
	if err := s.flushPartitions(); err != nil {
		return fmt.Errorf("failed to close storage: %w", err)
	}
	if err := s.removeExpiredPartitions(); err != nil {
		return fmt.Errorf("failed to remove expired partitions: %w", err)
	}
	return nil
}

func (s *Storage) removeExpiredPartitions() error {
	expiredList := make([]partition, 0)
	iterator := s.partitionList.newIterator()
	for iterator.next() {
		part := iterator.value()
		if part == nil {
			return fmt.Errorf("unexpected nil partition found")
		}
		if part.expired() {
			expiredList = append(expiredList, part)
		}
	}

	for i := range expiredList {
		if err := s.partitionList.remove(expiredList[i]); err != nil {
			return fmt.Errorf("failed to remove expired partition")
		}
	}
	return nil
}
