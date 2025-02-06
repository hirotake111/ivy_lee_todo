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

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: make([]*domain.Task, 0, defaultMemorySize),
	}
}

// Create implements domain.TaskRepository.
func (m *MemoryRepository) Create(ctx context.Context, db *db.Db, req *domain.NewTaskRequest) error {
	if len(req.Title) == 0 {
		return errors.New("title can't be empty")
	}
	id := len(m.tasks) + 1
	t := req.ToActionableTask(id)
	m.tasks = append(m.tasks, t)
	return nil
}

// Delete implements domain.TaskRepository.
func (m *MemoryRepository) Delete(ctx context.Context, db *db.Db, id int) error {
	i := id - 1
	if i < 0 || i >= m.len() {
		return fmt.Errorf("invalid ID '%d' specified", id)
	}
	m.tasks[i] = m.tasks[i].ToDeleted()
	return nil
}

// Find implements domain.TaskRepository.
func (m *MemoryRepository) Find(ctx context.Context, db *db.Db, id int) (*domain.Task, error) {
	for _, t := range m.tasks {
		if t.Id() == id {
			return t, nil
		}
	}
	return nil, apperrors.NotFound
}

// List implements domain.TaskRepository.
func (m *MemoryRepository) List(ctx context.Context, db *db.Db) ([]*domain.Task, error) {
	return m.tasks, nil
}

// Update implements domain.TaskRepository.
func (m *MemoryRepository) Update(ctx context.Context, db *db.Db, task *domain.Task) error {
	panic("unimplemented")
}

func (m *MemoryRepository) len() int {
	return len(m.tasks)
}

func (m MemoryRepository) debug() {
	for _, v := range m.tasks {
		fmt.Printf("actionable: %t,\tdeleted: %t\t%s\n", v.IsActionable(), v.IsDeleted(), v.Title())
	}
	fmt.Println("")
}
