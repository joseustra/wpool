package wpool

import (
	"sync"
)

// Job defines the task to be executed and how many times it will retry
// on fail
type Job struct {
	Task  func() error
	Retry int
}

// Pool defines the pool
type Pool struct {
	concurrency int
	queue       chan Job
	wg          sync.WaitGroup
}

// New prepares the pool and workers and size of total queue
func New(n, queueSize int) *Pool {
	pool := &Pool{
		concurrency: n,
		queue:       make(chan Job, queueSize),
	}

	for i := 0; i < n; i++ {
		w := newWorker(i+1, pool, &pool.wg)
		w.start()
	}

	return pool
}

// Add adds a new job to the queue
func (p *Pool) Add(j Job) {
	p.wg.Add(1)
	p.queue <- j
}

// Wait will waits for all tasks to be done before continue
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Drain will drain worker pool
func (p *Pool) Drain() {
	close(p.queue)
}

type worker struct {
	id   int
	pool *Pool
	wg   *sync.WaitGroup
}

func newWorker(id int, pool *Pool, wg *sync.WaitGroup) *worker {
	return &worker{
		id:   id,
		pool: pool,
		wg:   wg,
	}
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case j, ok := <-w.pool.queue:
				// if ok is false, it means that queue channel is closed
				if !ok {
					return
				}
				if j.Retry >= 0 {
					err := j.Task()
					if err != nil {
						j.Retry--
						w.pool.Add(j)
					}
				}
				w.wg.Done()
			}
		}
	}()
}
