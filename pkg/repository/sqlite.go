package repository

import (
	"context"

	"github.com/hirotake111/ivy_lee_todo/pkg/apperrors"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"

	_ "github.com/mattn/go-sqlite3"
)

type taskDto struct {
	Id          int
	Title       string
	Description string
}

func (t taskDto) toTask() *domain.Task {
	return domain.NewTask(t.Id, t.Title, t.Description)
}

type SQLiteRepository struct{}

func NewSQLiteRepository() *SQLiteRepository {
	return &SQLiteRepository{}
}

// Create implements domain.TaskRepository.
func (s *SQLiteRepository) Create(ctx context.Context, db db.Executor, t *domain.NewTaskRequest) error {
	_, err := db.Exec(ctx, `INSERT INTO task (title, description) VALUES ($1, $2);`, t.Title, t.Description)
	return err
}

// Delete implements domain.TaskRepository.
func (s *SQLiteRepository) Delete(ctx context.Context, db db.Executor, id int) error {
	_, err := db.Exec(ctx, `DELETE FROM task WHERE id = $1;`, id)
	return err
}

// Find implements domain.TaskRepository.
func (s *SQLiteRepository) Find(ctx context.Context, db db.Queryer, id int) (*domain.Task, error) {
	row := db.QueryRow(ctx, "SELECT id, title, description FROM task WHERE id = $1 and deleted_at is null", id)
	var t taskDto
	if err := row.Scan(&t.Id, &t.Title, &t.Description); err != nil {
		return nil, apperrors.NotFound
	}
	return t.toTask(), nil
}

// ListActionable implements domain.TaskRepository.
func (s *SQLiteRepository) ListActionable(ctx context.Context, db db.Queryer) (domain.TaskList, error) {
	rows, err := db.Query(ctx, "SELECT id, title, description FROM task WHERE actionable = 1 and deleted_at is null")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var l domain.TaskList
	for rows.Next() {
		var t taskDto
		if err := rows.Scan(&t.Id, &t.Title, &t.Description); err != nil {
			return nil, err
		}
		l = append(l, t.toTask())
	}
	return l, nil
}

// ListNonactionable implements domain.TaskRepository.
func (s *SQLiteRepository) ListNonactionable(ctx context.Context, db db.Queryer) ([]*domain.Task, error) {
	rows, err := db.Query(ctx, "SELECT id, title, description FROM task WHERE actionable = 0 and deleted_at is null")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var l domain.TaskList
	for rows.Next() {
		var t taskDto
		if err := rows.Scan(&t.Id, &t.Title, &t.Description); err != nil {
			return nil, err
		}
		l = append(l, t.toTask())
	}
	return l, nil
}

// Update implements domain.TaskRepository.
func (s *SQLiteRepository) Update(ctx context.Context, db db.Executor, t *domain.Task) error {
	_, err := db.Exec(
		ctx,
		"UPDATE task SET title = $1, description = $2, actionable = $3 WHERE id = $4",
		t.Title(), t.Description(), t.IsActionable(), t.Id(),
	)
	return err
}
