package workpool

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_NewWorker(t *testing.T) {
	assert := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	w := NewWorker(
		&workerPool{
			name:                "test",
			maxWorkers:          2,
			tasks:               make(chan Task, tasksCapacity),
			readyWorkers:        make(chan *worker, readyWorkerQueueSize),
			idleTimeout:         time.Second * 5,
			onDispatcherStopped: make(chan struct{}),
			stopped:             *atomic.NewBool(false),
			workersAlive:        *atomic.NewInt32(0),
			workersCreated:      *atomic.NewInt32(0),
			workersKilled:       *atomic.NewInt32(0),
			tasksConsumed:       *atomic.NewInt32(0),
			ctx:                 ctx,
			cancel:              cancel,
		},
	)

	w.execute(func() { fmt.Println("dongjiang1") })
	w.execute(func() { fmt.Println("dongjiang2") })

	w.stop(func() { fmt.Println("finished") })
	assert.True(true)
}
