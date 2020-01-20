// Package journal connects the user to the journal database and allows
// the user to alter the table journal_entries.
package journal

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var dbid int       // id value of table journal_entries
var dbdate string  // date value of table journal_entries
var dbentry string // entry value of table journal_entries

// InputEntry adds the current date as a string and prompts the user for
// a journal entry input to be stored into the database in association
// with the date.
func InputEntry(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(`Input journal entry:`)
	journalEntry, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	journalEntry = journalEntry[:len(journalEntry)-1]

	journalDate := string(time.Now().Format("01-02-2006"))

	ifEntryExists(db, journalEntry, journalDate)

	printEntry(db, journalDate)

}

// ViewEntry prints the date and entry of a particular date
func ViewEntry(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to view (MM-DD-YYYY):")
	journalDate, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
	journalDate = journalDate[:len(journalDate)-1]

	rows, err := db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&dbid, &dbdate, &dbentry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dbdate + ":\n" + dbentry)
	}
}

// ViewEntireJournal prints every date and entry of journal_entries
func ViewEntireJournal(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM journal_entries ORDER BY date")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&dbid, &dbdate, &dbentry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dbdate + ":\n" + dbentry)
	}
}

// DeleteEntry deletes the record of a particular date
func DeleteEntry(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to delete (MM-DD-YYYY):")
	journalDate, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	journalDate = journalDate[:len(journalDate)-1]

	statement, err := db.Prepare("DELETE FROM journal_entries WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	statement.Exec(journalDate)
}

// DeleteJournal deletes the entire table of journal_entries
func DeleteJournal(db *sql.DB) {
	statement, err := db.Prepare("DROP TABLE journal_entries")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

// EditEntry replaces the entry of a particular date
func EditEntry(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to edit:")
	journalDate, _ := reader.ReadString('\n')
	journalDate = journalDate[:len(journalDate)-1]

	rows, err := db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&dbid, &dbdate, &dbentry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dbdate + ":\n" + dbentry)
	}

	fmt.Println("Input replacement entry:")
	journalEntry, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	journalEntry = journalEntry[:len(journalEntry)-1]

	statement, err := db.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	statement.Exec(journalEntry, journalDate)

	printEntry(db, journalDate)
}

// InputEntryDate prompts the user for a date as a string and prompts the user
// for a journal entry input to be stored into the database in association
// with the date.
func InputEntryDate(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date (MM-DD-YYYY): ")
	journalDate, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	journalDate = journalDate[:len(journalDate)-1]

	fmt.Println("Input journal entry:")
	journalEntry, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	journalEntry = journalEntry[:len(journalEntry)-1]

	ifEntryExists(db, journalEntry, journalDate)

	printEntry(db, journalDate)
}

// ifEntryExists checks to see if an entry for a certain date already exists
// in the table journal_entires in a specified SQL database.
func ifEntryExists(db *sql.DB, journalEntry string, journalDate string) {
	rows, err := db.Query(`SELECT * FROM journal_entries WHERE date = ?`, journalDate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dateExists := false

	for rows.Next() {
		err := rows.Scan(&dbid, &dbdate, &dbentry)
		if err != nil {
			log.Fatal(err)
		}
		if journalDate == dbdate {
			dateExists = true
		}
	}

	// If the date of the entry already exists, the entry will be added	to
	// the preexisting entry after a new line.
	if dateExists {
		rows, err = db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			err := rows.Scan(&dbid, &dbdate, &dbentry)
			if err != nil {
				log.Fatal(err)
			}
			journalEntry = fmt.Sprint(dbentry + "\n\n" + journalEntry)
		}

		statement, err := db.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()
		statement.Exec(journalEntry, journalDate)

	} else {
		statement, err := db.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()
		statement.Exec(journalDate, journalEntry)
	}
}

// printEntry prints the entry of a specified date onto the console
func printEntry(db *sql.DB, journalDate string) {
	rows, err := db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&dbid, &dbdate, &dbentry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dbdate + ":\n" + dbentry)
	}
}
