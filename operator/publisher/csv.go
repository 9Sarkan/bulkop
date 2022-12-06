package publisher

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

type csvPublisher struct {
	fileReader *csv.Reader
	options    *Options
}

// NilReader if given csv reader are nil this error must be returned
var NilReader error = fmt.Errorf("reader shouldn't be nil")

// NewCsvPublisher create a new csv publisher with given options
// default option all flags have disabled
func NewCsvPublisher(reader *csv.Reader, options *Options) (*csvPublisher, error) {
	if reader == nil {
		return nil, NilReader
	}
	if options == nil {
		options = NewOptions(false, false, false)
	}
	return &csvPublisher{
		fileReader: reader,
		options:    options,
	}, nil
}

// Start read and publish on data channel
func (receiver csvPublisher) Start(ctx context.Context, dataChan chan interface{}) (int, error) {
	failure := 0
	for {
		select {
		case <-ctx.Done():
			return 0, nil
		default:
		}
		line, err := receiver.fileReader.Read()
		if err != nil {
			if err == io.EOF {
				return failure, nil
			}
			if err := receiver.handleFailure(err); err != nil {
				return receiver.options.failure, err
			}
		}
		// publish line in to channel
		dataChan <- line
	}
	return failure, nil
}

// handleFailure handle error if any error raised at read step
func (receiver csvPublisher) handleFailure(err error) error {
	if receiver.options.logFailure {
		log.Default().Println(err.Error())
	}
	if receiver.options.countFailure {
		receiver.options.failure++
	}
	if !receiver.options.continueOnErr {
		return err
	}
	return nil
}
