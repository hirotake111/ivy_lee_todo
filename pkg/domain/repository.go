package domain

import (
	"context"

	"github.com/hirotake111/ivy_lee_todo/pkg/db"
)

type NewTaskRequest struct {
	Title       string
	Description string
}

func (t *NewTaskRequest) ToActionableTask(id int) *Task {
	return &Task{
		id:          id,
		title:       t.Title,
		description: t.Description,
		actionable:  false,
	}
}

type TaskRepository interface {
	Create(ctx context.Context, db db.Executor, task *NewTaskRequest) error
	ListActionable(ctx context.Context, db db.Queryer) (TaskList, error)
	ListNonactionable(ctx context.Context, db db.Queryer) ([]*Task, error)
	Find(ctx context.Context, db db.Queryer, id int) (*Task, error)
	Update(ctx context.Context, db db.Executor, task *Task) error
	Delete(ctx context.Context, db db.Executor, id int) error
}
