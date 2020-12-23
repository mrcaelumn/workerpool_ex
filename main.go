package main

import (
	"fmt"
	"github.com/mrcaelumn/workerpool_ex/services"
)

func main() {

	maxCollectorWorker := 100

	collectorPool := services.NewWorkerPool(maxCollectorWorker)

	// start the worker
	collectorPool.StartWorkers()

	for i := 1; i < 100 ; i++ {
		data := fmt.Sprintf("this data %v", i)
		collectorPool.AddJob(data)
	}

}
