package workpool

// 执行任务的worker
type worker struct {
	pool   *workerPool
	tasks  chan Task
	stopCh chan struct{}
}

func NewWorker(pool *workerPool) *worker {
	w := &worker{
		pool:   pool,
		tasks:  make(chan Task),
		stopCh: make(chan struct{}),
	}
	w.pool.workersAlive.Inc()
	w.pool.workersCreated.Inc()
	go w.process()
	return w
}

func (w *worker) execute(task Task) {
	w.tasks <- task
}

func (w *worker) stop(callable func()) {
	defer callable()
	w.stopCh <- struct{}{}
	w.pool.workersKilled.Inc()
	w.pool.workersAlive.Dec()
}

func (w *worker) process() {
	var task Task
	for {
		select {
		case <-w.stopCh:
			return
		case task = <-w.tasks:
			task()
			w.pool.tasksConsumed.Inc()
			// 将w注册到readyWorkers
			w.pool.readyWorkers <- w
		}
	}
}
