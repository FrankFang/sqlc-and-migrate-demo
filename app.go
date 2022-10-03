package main

import (
	"context"
	"database/sql"
	"log"
	"reflect"

	"tutorial.sqlc.dev/app/tutorial"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

func run() error {
	ctx := context.Background()

	// 可以改数据库地址吗？
	db, err := sql.Open("sqlite3", "xxx.db")
	if err != nil {
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))

	// update the author's name to frank
	updatedAuthor, err := queries.UpdateAuthor(ctx, tutorial.UpdateAuthorParams{
		ID:   fetchedAuthor.ID,
		Name: "Frank",
		Bio:  sql.NullString{String: "", Valid: true},
	})

	if err != nil {
		return err
	}

	// prints true
	log.Println(updatedAuthor)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
