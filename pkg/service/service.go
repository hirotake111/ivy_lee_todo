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

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, s.db, id)
}

func (s *Service) ListActionableTask(ctx context.Context) (l []*domain.Task, err error) {
	ts, err := s.repo.List(ctx, s.db)
	if err != nil {
		return
	}
	for _, t := range ts {
		if !t.IsDeleted() && t.IsActionable() {
			l = append(l, t)
		}
	}
	return
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
