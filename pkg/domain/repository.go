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
	Create(ctx context.Context, db *db.Db, task *NewTaskRequest) error
	ListActionable(ctx context.Context, db *db.Db) (TaskList, error)
	ListNonactionable(ctx context.Context, db *db.Db) ([]*Task, error)
	Find(ctx context.Context, db *db.Db, id int) (*Task, error)
	Update(ctx context.Context, db *db.Db, task *Task) error
	Delete(ctx context.Context, db *db.Db, id int) error
}
