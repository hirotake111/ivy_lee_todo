package tui

import "github.com/hirotake111/ivy_lee_todo/pkg/domain"

type listMsg struct {
	tasks domain.TaskList
}

// errMsg is a message that contains an error in it.
type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

// newDisplayItemsMsg contains a new list of items currently been displayed
type newDisplayItemsMsg struct {
	tasks []*domain.Task
}

// successMsg contains text message indicating the result of an action
type successMsg struct {
	text string
}

// tasksUpdatedMsg indicates tasks are stable and UI need to reflesh
type tasksUpdatedMsg struct{}
