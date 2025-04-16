package domain

import "time"

type Task struct {
	id          int
	title       string
	description string
	actionable  bool
	createdAt   time.Time
}

func NewTask(id int, title, description string, actionable bool, now time.Time) *Task {
	return &Task{
		id:          id,
		title:       title,
		description: description,
		actionable:  actionable,
		createdAt:   now,
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

func (t Task) CreatedAt() time.Time {
	return t.createdAt
}

// ToActionable makes the task itself actionable
func (t *Task) ToActionable() *Task {
	t.actionable = true
	return t
}

// ToPlanned turns the task itself into planned one
func (t *Task) ToPlanned() *Task {
	t.actionable = false
	return t
}
