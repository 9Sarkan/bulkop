package operator

import "time"

// NewPipe create a pipe to handle request chan objects and put result on response channel
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
