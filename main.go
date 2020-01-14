package main

import (
	"database/sql"
	"flag"

	"github.com/ajkim19/project-0/journal"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Opens database connection
	database, _ := sql.Open("sqlite3", "./journal.db")

	// Creates table if it does not exist
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	statement.Exec()

	var view bool
	var delete bool
	var edit bool
	var all bool

	flag.BoolVar(&view, "view", false, "view entry")
	flag.BoolVar(&delete, "delete", false, "delete entry")
	flag.BoolVar(&edit, "edit", false, "edit entry")
	flag.BoolVar(&all, "all", false, "apply to every entry")
	flag.Parse()

	switch {
	case view == true:
		journal.ViewEntry(database)
	case view == true && all == true:
		journal.ViewEntireJournal(database)
	case delete == true:
		journal.DeleteEntry(database)
	case delete == true && all == true:
		journal.DeleteTable(database)
	case edit == true:
		journal.EditEntry(database)
	default:
		journal.InputEntry(database)
	}

}
