package domain

type Task struct {
	id          int
	title       string
	description string
	actionable  bool
	deleted     bool
}

func (t Task) Id() int {
	return t.id
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t Task) IsActionable() bool {
	return t.actionable
}

func (t *Task) ToActionable() *Task {
	t.actionable = true
	return t
}

func (t Task) IsDeleted() bool {
	return t.deleted
}

func (t *Task) ToDeleted() *Task {
	t.deleted = true
	return t
}
