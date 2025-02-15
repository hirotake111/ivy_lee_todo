package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hirotake111/ivy_lee_todo/pkg/apperrors"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
)

const (
	defaultMemorySize = 64
)

type MemoryRepository struct {
	tasks []*domain.Task
}

// ListNonactionable implements domain.TaskRepository.
func (m *MemoryRepository) ListNonactionable(ctx context.Context, db db.Queryer) (ts []*domain.Task, err error) {
	for _, t := range m.tasks {
		if !t.IsActionable() {
			ts = append(ts, t)
		}
	}
	return
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: make([]*domain.Task, 0, defaultMemorySize),
	}
}

// Create implements domain.TaskRepository.
func (m *MemoryRepository) Create(ctx context.Context, db db.Transaction, req *domain.NewTaskRequest) error {
	if len(req.Title) == 0 {
		return errors.New("title can't be empty")
	}
	id := len(m.tasks) + 1
	t := req.ToActionableTask(id)
	m.tasks = append(m.tasks, t)
	return nil
}

// Delete implements domain.TaskRepository.
func (m *MemoryRepository) Delete(ctx context.Context, db db.Transaction, id int) error {
	var tasks []*domain.Task
	for _, t := range m.tasks {
		if t.Id() != id {
			tasks = append(tasks, t)
		}
	}
	m.tasks = tasks
	return nil
}

// Find implements domain.TaskRepository.
func (m *MemoryRepository) Find(ctx context.Context, db db.Queryer, id int) (*domain.Task, error) {
	for _, t := range m.tasks {
		if t.Id() == id {
			return t, nil
		}
	}
	return nil, apperrors.NotFound
}

// List implements domain.TaskRepository.
func (m *MemoryRepository) ListActionable(ctx context.Context, db db.Queryer) (tl domain.TaskList, err error) {
	for _, t := range m.tasks {
		if t.IsActionable() {
			tl = append(tl, t)
		}
	}
	return
}

// Update implements domain.TaskRepository.
func (m *MemoryRepository) Update(ctx context.Context, db db.Transaction, task *domain.Task) error {
	for i, t := range m.tasks {
		if t.Id() == task.Id() {
			m.tasks[i] = task
		}
	}
	return nil
}

func (m MemoryRepository) debug() {
	for _, v := range m.tasks {
		fmt.Printf("actionable: %t,\t%s\n", v.IsActionable(), v.Title())
	}
	fmt.Println("")
}
