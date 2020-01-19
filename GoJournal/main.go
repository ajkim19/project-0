package main

import (
	"fmt"

	"github.com/ajkim19/project-0/GoJournal/journal"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Closes the database once the program has finished
	defer database.Close()

	fmt.Println("Welcome to GoJournal!")

	switch {
	case d == true:
		journal.InputEntryDate(database)
	case view == true && all == true:
		journal.ViewEntireJournal(database)
	case view == true:
		journal.ViewEntry(database)
	case delete == true && all == true:
		journal.DeleteJournal(database)
	case delete == true:
		journal.DeleteEntry(database)
	case edit == true:
		journal.EditEntry(database)
	default:
		journal.InputEntry(database)
	}
}
