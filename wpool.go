package wpool

import (
	"sync"
)

// Job is the definition of the task that will be executed
type Job func() error

// Pool defines the pool
type Pool struct {
	concurrency int
	queue       chan Job
	wg          sync.WaitGroup
}

// New prepares the pool and workers
func New(n int) *Pool {
	pool := &Pool{
		concurrency: n,
		queue:       make(chan Job, 100),
	}

	for i := 0; i < n; i++ {
		w := newWorker(i+1, pool.queue, &pool.wg)
		w.start()
	}

	return pool
}

// Add adds a new job to the queue
func (p *Pool) Add(j Job) error {
	p.wg.Add(1)
	p.queue <- j
	return nil
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
	id    int
	queue chan Job
	wg    *sync.WaitGroup
}

func newWorker(id int, queue chan Job, wg *sync.WaitGroup) *worker {
	return &worker{
		id:    id,
		queue: queue,
		wg:    wg,
	}
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case wkr, ok := <-w.queue:
				// if ok is false, it means that queue channel is closed
				if !ok {
					return
				}
				wkr()
				w.wg.Done()
			}
		}
	}()
}
