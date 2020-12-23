package services

import (
	"fmt"
)

// Job handles Get Request Payload
type Job struct {
	Payload string
}

// Worker represents the worker that executes the job
type Worker struct {
	// GET Job
	Pool       chan chan Job
	JobChannel chan Job

	quit chan bool
}

// NewWorker initiate worker
func NewWorker(pool chan chan Job) Worker {
	return Worker{
		Pool:       pool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.Pool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work from GET request.
				// do task
				fmt.Println(job.Payload)
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()

}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
