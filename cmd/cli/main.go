package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hirotake111/ivy_lee_todo/pkg/cli"
	"github.com/hirotake111/ivy_lee_todo/pkg/db"
	"github.com/hirotake111/ivy_lee_todo/pkg/repository"
)

func main() {
	db := db.NewDb()
	r := repository.NewMemoryRepository()
	c := cli.New(db, r)
	if err := c.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}
