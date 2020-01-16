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
var err1 error

func init() {
	database, err1 = sql.Open("sqlite3", "./journal.db")
	if err1 != nil {
		log.Fatal(err1)
	}

	statement, err2 := database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	if err2 != nil {
		log.Fatal(err2)
	}
	statement.Exec()

	flag.BoolVar(&view, "view", false, "view entry")
	flag.BoolVar(&delete, "delete", false, "delete entry")
	flag.BoolVar(&edit, "edit", false, "edit entry")
	flag.BoolVar(&all, "all", false, "apply to every entry")
	flag.Parse()
}
