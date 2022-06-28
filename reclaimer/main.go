package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sys/unix"
)

func main() {
	// Initialize SQLite
	log.Print("Initializing database...")
	const file = "/opt/kob/obituaries.db"
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err.Error())
	}
	log.Print("Success!")

	var stat unix.Statfs_t
	unix.Statfs("/opt/kob", &stat)

	// Available blocks * size per block = available space in bytes
	diskRemaining := float64(stat.Bavail*uint64(stat.Bsize)) / float64(stat.Blocks*uint64(stat.Bsize))

	for {
		err = reclaim(diskRemaining, db)
		if err != nil {
			log.Print("Cannot get row count, error: ", err.Error())
		}
		time.Sleep(60 * time.Second)
	}
}

// If there's less than 10% of disk remaining, delete 10% of rows and reclaim 10% of disk..
func reclaim(diskRemaining float64, db *sql.DB) error {
	var err error
	if diskRemaining < .1 {
		_, err := db.Exec("DELETE FROM pods WHERE rowid IN (SELECT rowid FROM pods ORDER BY rowid LIMIT (SELECT COUNT(*)/10 FROM pods));")
		if err != nil {
			return err
		}
		log.Print("Deleted 10% of entries.")
	}
	return err
}
