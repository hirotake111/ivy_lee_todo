package service

import (
	"context"

	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
)

type Service struct {
	db   *db.Db
	repo domain.TaskRepository
}

// TODO: implement
func (s *Service) ListTasks(ctx context.Context) ([]domain.Task, error) {
	return s.repo.List(ctx, s.db)
}

func (s *Service) AddTask(ctx context.Context, title, description string) error {
	req := domain.NewTaskRequest{
		Title:       title,
		Description: description,
	}
	return s.repo.Create(ctx, s.db, &req)
}

func NewService(db *db.Db, repo domain.TaskRepository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}
