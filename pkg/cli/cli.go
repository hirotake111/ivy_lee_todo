package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
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
LOOP:
	for {
		if err := c.list(ctx); err != nil {
			fmt.Printf("Error listing tasks: %s\n", err)
			break
		}
		cmd, err := c.getCommand()
		if err != nil {
			fmt.Printf("Error getting a command: %s\n", err)
			break
		}
		switch cmd {
		case "a":
			if err := c.add(ctx); err != nil {
				fmt.Printf("Error adding a new task: %s\n", err)
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

func (c Cli) getCommand() (string, error) {
	fmt.Printf("\n[a]dd a new task\t[q]uit program\nEnter command:")
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
	return c.service.AddTask(ctx, title, "")
}

// list task mode
func (c *Cli) list(ctx context.Context) error {
	tasks, err := c.service.ListTasks(ctx)
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
