package worker

import "sync"

type Pool struct {
	taskQueue   chan func()
	waitGroup   sync.WaitGroup
	concurrency int
}

func NewWorkerPool(concurrency int) *Pool {
	return &Pool{
		taskQueue:   make(chan func()),
		concurrency: concurrency,
	}
}

func (wp *Pool) Start() {
	for i := 0; i < wp.concurrency; i++ {
		wp.waitGroup.Add(1)
		go wp.worker()
	}
}

func (wp *Pool) worker() {
	defer wp.waitGroup.Done()
	for task := range wp.taskQueue {
		task()
	}
}

func (wp *Pool) Submit(task func()) {
	wp.taskQueue <- task
}

func (wp *Pool) Stop() {
	close(wp.taskQueue)
	wp.waitGroup.Wait()
}
