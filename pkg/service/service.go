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
	tl, err := s.repo.ListActionable(ctx, s.db)
	if err != nil {
		return err
	}
	if !tl.CanAddAnother() {
		return apperrors.NewTaskExceededError(tl.MaxTskNum())
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	t, err := s.repo.Find(ctx, tx, id)
	if err != nil {
		return err
	}
	t.ToActionable()
	if err := s.repo.Update(ctx, tx, t); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) Update(ctx context.Context, t *domain.Task) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.repo.Update(ctx, tx, t); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) Find(ctx context.Context, id int) (*domain.Task, error) {
	return s.repo.Find(ctx, s.db, id)
}

func (s *Service) ListPlannedTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.ListNonactionable(ctx, s.db)
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) ListActionableTask(ctx context.Context) (domain.TaskList, error) {
	return s.repo.ListActionable(ctx, s.db)
}

func (s *Service) AddTask(ctx context.Context, title, description string) error {
	req := domain.NewTaskRequest{
		Title:       title,
		Description: description,
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.repo.Create(ctx, tx, &req); err != nil {
		return err
	}
	return tx.Commit()
}
