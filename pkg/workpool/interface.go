package workpool

import "time"

// TODO Task 应该有输入和输出
type Task func()

// 开辟一个协程池：当有任务提交时，提交到协程池中运行；如果协程池都在工作，任务挂起
type Pool interface {
	// 提交任务
	Submit(task Task)        // 提交任务
	SubmitAndWait(task Task) // 提交任务并等待其执行
	Stopped() bool           // 如果协程停止，返回true
	Stop()                   // 停下来优雅地停止所有的勾当，所有挂起的任务将在退出前完成
	//metrics
}

func NewDefaultPool(name string, maxWorkers int, idleTimeout time.Duration) Pool {
	return NewWorkerPool(name, maxWorkers, idleTimeout)
}
