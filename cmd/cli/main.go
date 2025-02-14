package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hirotake111/ivy_lee_todo/pkg/cli"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/repository"
	"github.com/hirotake111/ivy_lee_todo/pkg/service"
)

func main() {
	db := db.NewSqlite3Db(false)
	r := repository.NewSQLiteRepository()
	// r := repository.NewMemoryRepository()
	s := service.NewService(db, r)
	c := cli.New(s)
	if err := c.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}
