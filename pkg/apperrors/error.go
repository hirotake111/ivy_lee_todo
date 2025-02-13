package apperrors

import "errors"

type TaskNumbersExceededError struct {
	s          string
	maxTaskNum int
}

var (
	NotFound error = errors.New("not found")
	Quit     error = errors.New("quit")
)

func NewTaskExceededError(maxTaskNum int) TaskNumbersExceededError {
	return TaskNumbersExceededError{
		s:          "max task numbers exceeded",
		maxTaskNum: maxTaskNum,
	}
}

func (e TaskNumbersExceededError) Error() string {
	return e.s
}

func (e TaskNumbersExceededError) MaxTaskNum() int {
	return e.maxTaskNum
}
