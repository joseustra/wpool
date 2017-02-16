package wpool

import "fmt"

var workerQueue chan chan Job
var workQueue chan Job

func init() {
	workQueue = make(chan Job, 10)
}

// Job the function to be executed
type Job func() error

// WPool Handle all background jobs
type WPool struct {
	Size int
}

// New start all workers and return WPool
func New(size int) *WPool {
	workerQueue = make(chan chan Job, size)
	wp := &WPool{Size: size}

	for i := 0; i < wp.Size; i++ {
		w := newWorker(i+1, workerQueue)
		w.start()
	}

	go func() {
		for {
			select {
			case work := <-workQueue:
				go func() {
					w := <-workerQueue
					w <- work
				}()
			}
		}
	}()

	fmt.Printf("[WPool] %v workers started\n", wp.Size)

	return wp
}

// Add adds a new job to the queue
func (wp *WPool) Add(job Job) {
	workQueue <- job
}

type worker struct {
	id    int
	work  chan Job
	queue chan chan Job
}

func newWorker(id int, queue chan chan Job) *worker {
	return &worker{
		id:    id,
		work:  make(chan Job),
		queue: queue,
	}
}

func (w *worker) start() {
	go func() {
		for {
			w.queue <- w.work

			select {
			case work := <-w.work:
				err := work()
				if err != nil {
					fmt.Printf("Job failed %v\n", w.id)
				}
			}
		}
	}()
}
