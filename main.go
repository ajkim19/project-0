package main

import (
	"github.com/ajkim19/project-0/journal"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	journal.ViewEntireJournal(database)
	// switch {
	// case view == true && all == true:
	// 	journal.ViewEntireJournal(database)
	// case view == true:
	// 	journal.ViewEntry(database)
	// case delete == true && all == true:
	// 	journal.DeleteJournal(database)
	// case delete == true:
	// 	journal.DeleteEntry(database)
	// case edit == true:
	// 	journal.EditEntry(database)
	// default:
	// 	journal.InputEntry(database)
	// }
}
