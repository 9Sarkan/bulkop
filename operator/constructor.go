package operator

import "time"

func NewLine() Line {
	return Line{}
}

// NewPipe create a pipe to handle request chan objects and put result on response channel
func NewPipe(name string, requestChan chan interface{}, processor Processor,
	worker int, errorHandler ErrorHandler, processorTimeout time.Duration) Pipe {
	pipe := Pipe{
		name:             name,
		processor:        processor,
		errorHandler:     errorHandler,
		requestChan:      requestChan,
		responseChan:     make(chan interface{}),
		centralKiller:    make(chan struct{}),
		processorTimeout: processorTimeout,
	}
	for i := 0; i < worker; i++ {
		pipe.NewWorker()
	}
	return pipe
}
