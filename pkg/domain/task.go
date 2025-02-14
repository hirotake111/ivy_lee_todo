package domain

type Task struct {
	id          int
	title       string
	description string
	actionable  bool
}

func NewTask(id int, title, description string) *Task {
	return &Task{
		id:          id,
		title:       title,
		description: description,
		actionable:  false,
	}
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

// ToActionable makes the task itself actionable
func (t *Task) ToActionable() *Task {
	t.actionable = true
	return t
}
