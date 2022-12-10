package operator

import (
	"context"
	"log"
	"time"
)

// Pipe a stage of process of an entity
// with specific error handler and also response channels
type Pipe struct {
	name             string
	processor        Processor
	workerCount      int
	errorHandler     ErrorHandler
	requestChan      chan interface{}
	responseChan     chan interface{}
	centralKiller    chan struct{}
	processorTimeout time.Duration
}

func NewPipe(name string, requestChan, responseChan chan interface{}, processor Processor,
	worker int, errorHandler ErrorHandler, processorTimeout time.Duration) {
	pipe := Pipe{
		name:             name,
		processor:        processor,
		errorHandler:     errorHandler,
		requestChan:      requestChan,
		responseChan:     responseChan,
		centralKiller:    make(chan struct{}),
		processorTimeout: processorTimeout,
	}
	for i := 0; i < worker; i++ {
		pipe.NewWorker()
	}
}

// NewWorker create new worker to listen in request channel and put response in response channel
func (receiver Pipe) NewWorker() {
	go func() {
		for {
			select {
			case <-receiver.centralKiller:
				log.Default().Printf("Pipe %s worker stop according to central killer event", receiver.name)
				receiver.workerCount--
				return
			case data := <-receiver.requestChan:
				ctx, cancel := context.WithTimeout(context.Background(), receiver.processorTimeout)
				result, err := receiver.processor(ctx, data)
				if err != nil {
					if err := receiver.errorHandler(data, err); err != nil {
						// if errorHandler function return any error that
						// means to break and stop the worker
						receiver.workerCount--
						return
					}
				}
				cancel()
				receiver.responseChan <- result
			}
		}
	}()
	receiver.workerCount++
	log.Default().Printf("Pipe %s worker count increase to %d", receiver.name, receiver.workerCount)
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
