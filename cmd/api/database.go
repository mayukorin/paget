package main

import (
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

func init() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("error opening database: %q\n", err)
	}
	fmt.Printf(db.Stats().WaitDuration.String())
}
