package services

import (
	"strconv"

	"github.com/inconshreveable/log15"
)

var (
	startedWorkers []Worker
)

// WorkerPool is worker pooler management
type WorkerPool struct {
	pool     chan chan Job
	jobQueue chan Job


	quit       chan bool
	maxWorkers int

}

// NewWorkerPool creates new WorkerPool instance
func NewWorkerPool(maxWorkers int) WorkerPool {
	// Get
	pool := make(chan chan Job, maxWorkers)
	jobQueue := make(chan Job)


	startedWorkers = make([]Worker, maxWorkers)

	return WorkerPool{
		pool:     pool,
		jobQueue: jobQueue,

		quit:       make(chan bool),
		maxWorkers: maxWorkers,
	}
}

// AddJob adding a GET job to queue
func (p WorkerPool) AddJob(payload string) {
	p.jobQueue <- Job{Payload: payload}
}


// StartWorkers start all workers
func (p WorkerPool) StartWorkers() {
	log15.Info("starting collector workers: " + strconv.Itoa(p.maxWorkers) + " workers")
	for i := 0; i < p.maxWorkers; i++ {
		worker := NewWorker(p.pool)
		worker.Start()

		// register started worker
		startedWorkers = append(startedWorkers, worker)
	}

	go p.dispatchJobQueue()
}

// StopWorkers stop all workers
func (p WorkerPool) StopWorkers() {
	log15.Info("stopping job dispatcher")
	go func() {
		p.quit <- true
	}()

	log15.Info("stopping collector workers")
	for _, worker := range startedWorkers {
		w := worker
		w.Stop()
	}
}

func (p WorkerPool) dispatchJobQueue() {
	for {
		select {
		case job := <-p.jobQueue:
			log15.Debug("job received")
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-p.pool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		case <-p.quit:
			// we have received a signal to stop
			return
		}
	}
}
