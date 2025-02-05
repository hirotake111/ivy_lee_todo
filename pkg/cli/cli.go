package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hirotake111/ivy_lee_todo/pkg/apperrors"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/domain"
)

type Cli struct {
	db     *db.Db
	repo   domain.TaskRepository
	reader *bufio.Reader
}

func New(db *db.Db, repo domain.TaskRepository) *Cli {
	reader := bufio.NewReader(os.Stdin)
	return &Cli{
		repo:   repo,
		db:     db,
		reader: reader,
	}
}

func (c *Cli) Run(ctx context.Context) error {
LOOP:
	for {
		if err := c.listTasks(ctx); err != nil {
			break
		}
		fmt.Printf("\n[a]dd a new task\t[q]uit program\nEnter command:")
		cmd, err := c.reader.ReadString('\n')
		if err != nil {
			return err
		}
		switch strings.TrimSpace(cmd) {
		case "a":
			fmt.Println("You'are adding a new task")
			if err := c.AddTask(ctx); err != nil {
				fmt.Printf("error: %s\n", err.Error())
			}
		case "q":
			fmt.Println("Quitting program...")
			break LOOP
		default:
			fmt.Println("")
		}
	}
	return nil
}

func (c *Cli) AddTask(ctx context.Context) error {
	fmt.Printf("\n\nEnter title:")
	title, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	req := domain.NewTaskRequest{
		Title:       strings.TrimSpace(title),
		Description: "", // TODO: implement
	}
	return c.repo.Create(ctx, c.db, &req)
}

// listTasks displays a listTasks of actionable tasks
func (c *Cli) listTasks(ctx context.Context) error {
	tasks, err := c.repo.List(ctx, c.db)
	if err != nil && !errors.Is(err, apperrors.NotFound) {
		return err
	}
	fmt.Printf("==== Tasks ====\n\n")
	if len(tasks) == 0 {
		fmt.Println("\nNo Tasks")
		return nil
	}
	for i, t := range tasks {
		fmt.Printf("%d - %s\n", i+1, t.Title())
	}
	return nil
}
