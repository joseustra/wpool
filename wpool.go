package wpool

import "sync"

// Job is the definition of the task that will be executed
type Job func() error

// Pool defines the pool
type Pool struct {
	concurrency int
	queue       chan Job
}

var wg sync.WaitGroup

// New prepares the pool and workers
func New(n int) *Pool {
	pool := &Pool{
		concurrency: n,
		queue:       make(chan Job, 100),
	}

	for i := 0; i < n; i++ {
		w := newWorker(i+1, pool.queue)
		w.start()
	}

	return pool
}

// Add adds a new job to the queue
func (p *Pool) Add(j Job) error {
	wg.Add(1)
	p.queue <- j
	return nil
}

// Wait will waits for all tasks to be done before continue
func (p *Pool) Wait() {
	wg.Wait()
}

type worker struct {
	id    int
	queue chan Job
}

func newWorker(id int, queue chan Job) *worker {
	return &worker{id: id, queue: queue}
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case wkr := <-w.queue:
				wkr()
				wg.Done()
			}
		}
	}()
}
