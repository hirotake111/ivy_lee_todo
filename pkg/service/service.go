package service

import (
	"context"
	"fmt"
	"log"

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
	t, err := s.repo.Find(ctx, s.db, id)
	log.Printf("debug task: %+v\n", t)
	if err != nil {
		return err
	}
	t.ToActionable()
	return s.repo.Update(ctx, s.db, t)
}

func (s *Service) Update(ctx context.Context, t *domain.Task) error {
	return s.repo.Update(ctx, s.db, t)
}

func (s *Service) Find(ctx context.Context, id int) (*domain.Task, error) {
	return s.repo.Find(ctx, s.db, id)
}

func (s *Service) ListPlannedTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.ListNonactionable(ctx, s.db)
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	fmt.Println("deletetask ()")
	return s.repo.Delete(ctx, s.db, id)
}

func (s *Service) ListActionableTask(ctx context.Context) (domain.TaskList, error) {
	return s.repo.ListActionable(ctx, s.db)
}

func (s *Service) AddTask(ctx context.Context, title, description string) error {
	req := domain.NewTaskRequest{
		Title:       title,
		Description: description,
	}
	return s.repo.Create(ctx, s.db, &req)
}
