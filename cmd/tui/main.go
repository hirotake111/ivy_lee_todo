package main

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/repository"
	"github.com/hirotake111/ivy_lee_todo/pkg/service"
	"github.com/hirotake111/ivy_lee_todo/pkg/tui"
)

func main() {
	ctx := context.Background()
	db := db.NewSqlite3Db(false)
	r := repository.NewSQLiteRepository()
	s := service.NewService(db, r)
	m := tui.InitializeModel(ctx, s)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Couldn't start program: %s\n", err)
	}
}
