package operator

import "context"

type Publisher interface {
	Start(context.Context, chan interface{}) (int, error)
}

type Processor func(context.Context, interface{}) (interface{}, error)
type ErrorHandler func(interface{}, error) error
