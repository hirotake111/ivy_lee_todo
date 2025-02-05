package repository

import (
	"context"

	"github.com/hirotake111/ivy_lee_todo/pkg/apperrors"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
)

type MemoryRepository struct {
	tasks []domain.Task
}

// Create implements domain.TaskRepository.
func (m *MemoryRepository) Create(ctx context.Context, db *db.Db, req *domain.NewTaskRequest) error {
	id := len(m.tasks)
	t := req.ToTask(id)
	m.tasks = append(m.tasks, *t)
	return nil
}

// Delete implements domain.TaskRepository.
func (m *MemoryRepository) Delete(ctx context.Context, db *db.Db, id int) error {
	panic("unimplemented")
}

// Find implements domain.TaskRepository.
func (m *MemoryRepository) Find(ctx context.Context, db *db.Db, id int) (domain.Task, error) {
	for _, t := range m.tasks {
		if t.Id() == id {
			return t, nil
		}
	}
	return domain.Task{}, apperrors.NotFound
}

// List implements domain.TaskRepository.
func (m *MemoryRepository) List(ctx context.Context, db *db.Db) ([]domain.Task, error) {
	return m.tasks, nil
}

// Update implements domain.TaskRepository.
func (m *MemoryRepository) Update(ctx context.Context, db *db.Db, task *domain.Task) error {
	panic("unimplemented")
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: []domain.Task{},
	}
}
