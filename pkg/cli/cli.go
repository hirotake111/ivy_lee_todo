package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hirotake111/ivy_lee_todo/pkg/apperrors"
	"github.com/hirotake111/ivy_lee_todo/pkg/service"
)

type Cli struct {
	service *service.Service
	reader  *bufio.Reader
}

func New(service *service.Service) *Cli {
	reader := bufio.NewReader(os.Stdin)
	return &Cli{
		service: service,
		reader:  reader,
	}
}

func (c *Cli) Run(ctx context.Context) error {
	var displayedPlanned bool
LOOP:
	for {
		if displayedPlanned {
			displayedPlanned = false
			if err := c.listPlannedTasks(ctx); err != nil {
				fmt.Printf("Error listing non-actionable tasks: %s\n", err)
			}
		} else {
			if err := c.list(ctx); err != nil {
				fmt.Printf("Error listing tasks: %s\n", err)
				break
			}
		}
		cmd, err := c.getCommand()
		if err != nil {
			fmt.Printf("Error getting a command: %s\n", err)
			break
		}
		switch cmd {
		case "l":
			displayedPlanned = true
		case "a":
			if err := c.add(ctx); err != nil {
				fmt.Printf("Error adding a new task: %s\n", err)
			}
		case "d":
			if err := c.delete(ctx); err != nil {
				fmt.Printf("\nError deleting a new task: %s\n", err)
			}
		case "m":
			if err := c.makeActionable(ctx); err != nil {
				var exceededErr apperrors.TaskNumbersExceededError
				if errors.As(err, &exceededErr) {
					fmt.Printf("\nYou can't have more than %d tasks", exceededErr.MaxTaskNum())
				} else {
					fmt.Printf("\nError making a task actionable: %s\n", err)
				}
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

func (c *Cli) listPlannedTasks(ctx context.Context) any {
	tasks, err := c.service.ListPlannedTasks(ctx)
	if err != nil && !errors.Is(err, apperrors.NotFound) {
		return err
	}
	fmt.Printf("\n\n==== Planned Tasks ====\n\n")
	if len(tasks) == 0 {
		fmt.Println("\nNo Actionable Tasks")
		return nil
	}
	for _, t := range tasks {
		fmt.Printf("%d - %s\n", t.Id(), t.Title())
	}
	return nil
}

// delete task mode
func (c *Cli) delete(ctx context.Context) error {
	idStr, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return c.service.DeleteTask(ctx, id)
}

func (c Cli) getCommand() (string, error) {
	fmt.Printf("\n[l]ist planned tasks [a]dd a new task  [d]elete a task [m]ove task to actionable list  [q]uit program\nEnter command:")
	cmd, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(cmd), nil
}

// Add task mode
func (c *Cli) add(ctx context.Context) error {
	fmt.Printf("\n\nEnter title:")
	title, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	title = strings.TrimSpace(title)
	return c.service.AddTask(ctx, title, "")
}

// make a task actionable
func (c *Cli) makeActionable(ctx context.Context) error {
	fmt.Printf("\n\nEnter ID:")
	idStr, err := c.reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return c.service.MakeActionable(ctx, id)
}

// list task mode
func (c *Cli) list(ctx context.Context) error {
	tasks, err := c.service.ListActionableTask(ctx)
	if err != nil && !errors.Is(err, apperrors.NotFound) {
		return err
	}
	fmt.Printf("\n\n==== Tasks ====\n\n")
	if len(tasks) == 0 {
		fmt.Println("\nNo Tasks")
		return nil
	}
	for _, t := range tasks {
		fmt.Printf("%d - %s\n", t.Id(), t.Title())
	}
	return nil
}
