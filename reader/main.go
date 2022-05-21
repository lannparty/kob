package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize SQLite
	log.Print("Opening database...")
	const file = "/opt/kube-obituaries/obituaries.db"
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Success!")

	db.Exec("")
}
