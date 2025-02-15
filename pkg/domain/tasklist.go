package domain

const (
	maxTaskNum = 6
)

type TaskList []*Task

func NewTaskList(in []*Task) TaskList {
	return in
}

func (tl TaskList) CanAddAnother() bool {
	return len(tl.ActionableTasks()) < maxTaskNum
}

func (tl TaskList) ActionableTasks() []*Task {
	var l []*Task
	for _, t := range tl {
		if t.IsActionable() {
			l = append(l, t)
		}
	}
	return l
}

func (tl TaskList) PlannedTasks() []*Task {
	var l []*Task
	for _, t := range tl {
		if !t.IsActionable() {
			l = append(l, t)
		}
	}
	return l
}

func (tl TaskList) MaxTskNum() int {
	return maxTaskNum
}
