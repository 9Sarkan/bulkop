package operator

import "context"

type Pipe struct {
	processor    Processor
	workerCount  int
	resultChan   chan interface{}
	errorHandler func(interface{}, error)
}

func (receiver Pipe) NewWorker(ctx context.Context, dataChan <-chan interface{}) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-dataChan:
				result, err := receiver.processor(ctx, data)
				if err != nil {
					receiver.errorHandler(data, err)
				}
				receiver.resultChan <- result
			}
		}
	}()
}
