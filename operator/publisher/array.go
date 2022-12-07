package publisher

import (
	"context"
	"fmt"
)

type arrayPublisher struct {
	data          []interface{}
	lastReadIndex int
}

// EmptyArrayErr given data array is empty
var EmptyArrayErr error = fmt.Errorf("array is empty")

// NewArrayPublisher create a new array publisher with data source of slice type
func NewArrayPublisher(data []interface{}) (*arrayPublisher, error) {
	if len(data) == 0 {
		return nil, EmptyArrayErr
	}
	return &arrayPublisher{
		data: data,
	}, nil
}

// Start read from array and push to data channel
func (receiver arrayPublisher) Start(ctx context.Context, dataChannel chan interface{}) (int, error) {
	startFrom := receiver.lastReadIndex
	for _, obj := range receiver.data[startFrom:] {
		select {
		case <-ctx.Done():
			return len(receiver.data) - (receiver.lastReadIndex + 1), nil
		default:
		}
		dataChannel <- obj
		receiver.lastReadIndex++
	}
	return 0, nil
}
