package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

var d bool           // Flag to input an entry for a particular date
var view bool        // Flag to view table journal_entries
var delete bool      // Flag to delete from table journal_entries
var edit bool        // Flag to edit journal_entries
var all bool         // Flag to apply alteration to the entire table of journal_entries
var flag1 string     // First flag string
var flag2 string     // Second flag string
var database *sql.DB // Pointer to database handle
var err error        // Temporary storage of error value
var id int           // Temporary storage of id value of table journal_entries
var date string      // Temporary storage of date value of table journal_entries
var entry string     // Temporary storage of entry value of table journal_entries

func init() {
	// Makes a handle for the database journal
	database, err = sql.Open("sqlite3", "./journal.db")
	if err != nil {
		log.Fatal(err)
	}

	// Creates the table journal_entries if it has been dropped
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	rows, err := database.Query("SELECT * FROM journal_entries")
	if err != nil {
		log.Fatal(err)
	}

	var dateExists bool

	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		if date == "01-01-2020" {
			dateExists = true
		}
	}

	// Adds an entry to journal_entries if it is empty
	if dateExists == false {
		statement, err := database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-01-2020", "Today is New Year's Day!")`)
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec()
	}

	// Initalizes the flags
	flag.BoolVar(&d, "d", false, "add entry to specified date")
	flag.BoolVar(&view, "view", false, "view entry")
	flag.BoolVar(&delete, "delete", false, "delete entry")
	flag.BoolVar(&edit, "edit", false, "edit entry")
	flag.BoolVar(&all, "all", false, "apply to every entry")
	flag.Parse()

	// Removes special characters of flags
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 2 {
		flag1 = reg.ReplaceAllString(os.Args[1], "")
	} else if len(os.Args) > 2 {
		flag1 = reg.ReplaceAllString(os.Args[1], "")
		flag2 = reg.ReplaceAllString(os.Args[2], "")
	}
}
