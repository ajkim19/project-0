package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var view bool
var delete bool
var edit bool
var all bool
var database *sql.DB
var err error
var id int
var date string
var entry string
var dateExists bool

func init() {
	database, err = sql.Open("sqlite3", "./journal.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	rows, err := database.Query("SELECT * FROM journal_entries")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		if date == "01-01-2020" {
			dateExists = true
		}
	}

	if dateExists == false {
		statement, err := database.Prepare("INSERT INTO journal_entries (date, entry) VALUES ('01-01-2020', 'Today is New Year's Day!')")
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec()
	}

	flag.BoolVar(&view, "view", false, "view entry")
	flag.BoolVar(&delete, "delete", false, "delete entry")
	flag.BoolVar(&edit, "edit", false, "edit entry")
	flag.BoolVar(&all, "all", false, "apply to every entry")
	flag.Parse()
}
