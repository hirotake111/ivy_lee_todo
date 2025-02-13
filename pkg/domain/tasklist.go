package domain

type TaskList []*Task

const (
	maxTaskNum = 6
)

func NewTaskList(in []*Task) TaskList {
	return in
}

func (t TaskList) CanAddAnother() bool {
	return len(t) < maxTaskNum
}

func (t TaskList) MaxTskNum() int {
	return maxTaskNum
}
