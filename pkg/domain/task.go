package domain

type Task struct {
	id           int
	title        string
	description  string
	isActionable bool
}

func NewTask(id int, title string, desc string) *Task {
	return &Task{
		id:           id,
		title:        title,
		description:  desc,
		isActionable: false,
	}
}

func (t *Task) Id() int {
	return t.id
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

func (t *Task) IsActionable() bool {
	return t.isActionable
}
