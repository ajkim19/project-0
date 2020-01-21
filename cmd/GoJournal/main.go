package main

import (
	"fmt"

	"github.com/ajkim19/project-0/pkg/journal"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Closes the database once the program has finished
	defer database.Close()

	fmt.Println(" _______________________")
	fmt.Println("|                       |")
	fmt.Println("| Welcome to GoJournal! |")
	fmt.Println("|_______________________|\n")

	switch flag1 {
	case "date":
		journal.InputEntryDate(database)
	case "view":
		if flag2 == "all" {
			journal.ViewEntireJournal(database)
		} else {
			journal.ViewEntry(database)
		}
	case "delete":
		if flag2 == "all" {
			journal.DeleteJournal(database)
		} else {
			journal.DeleteEntry(database)
		}
	case "edit":
		journal.EditEntry(database)
	default:
		journal.InputEntry(database)
	}
}
