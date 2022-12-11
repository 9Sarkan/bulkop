package operator

import "errors"

// Line for process entities through pipes
type Line struct {
	pipes       []Pipe
	finalResult chan interface{}
}

var EmptyLine error = errors.New("line is empty")

// AddPipe add a pipe to end of the line
func (receiver Line) AddPipe(pipe Pipe) {
	receiver.pipes = append(receiver.pipes, pipe)
}

// GetFinalResultChan return final result channel
func (receiver Line) GetFinalResultChan() (chan interface{}, error) {
	if len(receiver.pipes) == 0 {
		return nil, EmptyLine
	}
	return receiver.pipes[len(receiver.finalResult)-1].GetResponseChan(), nil
}
