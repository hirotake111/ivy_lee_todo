package service

import (
	"context"

	"github.com/hirotake111/ivy_lee_todo/pkg/apperrors"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
)

type Service struct {
	db   *db.Db
	repo domain.TaskRepository
}

func NewService(db *db.Db, repo domain.TaskRepository) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

func (s *Service) MakeActionable(ctx context.Context, id int) error {
	tl, err := s.repo.List(ctx, s.db)
	if err != nil {
		return err
	}
	if !tl.CanAddAnother() {
		return apperrors.NewTaskExceededError(tl.MaxTskNum())
	}

	return s.db.StartTransaction(ctx, func(tx db.Transaction) error {
		t, err := s.repo.Find(ctx, tx, id)
		if err != nil {
			return err
		}
		t.ToActionable()
		return s.repo.Update(ctx, tx, t)
	})

}

func (s *Service) Update(ctx context.Context, t *domain.Task) error {
	return s.db.StartTransaction(ctx, func(tx db.Transaction) error {
		return s.repo.Update(ctx, tx, t)
	})
}

func (s *Service) Find(ctx context.Context, id int) (*domain.Task, error) {
	return s.repo.Find(ctx, s.db, id)
}

func (s *Service) ListPlannedTasks(ctx context.Context) ([]*domain.Task, error) {
	l, err := s.repo.List(ctx, s.db)
	return l.PlannedTasks(), err
}

func (s *Service) ListActionableTask(ctx context.Context) (domain.TaskList, error) {
	l, err := s.repo.List(ctx, s.db)
	return l.ActionableTasks(), err
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	return s.db.StartTransaction(ctx, func(tx db.Transaction) error {
		return s.repo.Delete(ctx, tx, id)
	})
}

func (s *Service) AddTask(ctx context.Context, title, description string) error {
	req := domain.NewTaskRequest{
		Title:       title,
		Description: description,
	}
	return s.db.StartTransaction(ctx, func(tx db.Transaction) error {
		return s.repo.Create(ctx, tx, &req)
	})
}
