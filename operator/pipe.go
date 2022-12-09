package operator

import (
	"context"
	"log"
)

// Pipe a stage of process of an entity
// with specific error handler and also response channels
type Pipe struct {
	name          string
	processor     Processor
	workerCount   int
	errorHandler  func(interface{}, error)
	requestChan   chan interface{}
	responseChan  chan interface{}
	centralKiller chan struct{}
}

// NewWorker create new worker to listen in request channel and put response in response channel
func (receiver Pipe) NewWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Default().Printf("Pipe %s worker stop according to context", receiver.name)
				return
			case <-receiver.centralKiller:
				log.Default().Printf("Pipe %s worker stop according to central killer event", receiver.name)
				return
			case data := <-receiver.requestChan:
				result, err := receiver.processor(ctx, data)
				if err != nil {
					receiver.errorHandler(data, err)
				}
				receiver.responseChan <- result
			}
		}
	}()
	receiver.workerCount++
}

// DecreaseWorkerCount decrease one of workers
func (receiver Pipe) DecreaseWorkerCount() {
	// check if there is an active worker kill it
	if receiver.workerCount > 0 {
		receiver.centralKiller <- struct{}{}
		receiver.workerCount--
	}
}

// GetWorkerCount worker count getter
func (receiver Pipe) GetWorkerCount() int {
	return receiver.workerCount
}
