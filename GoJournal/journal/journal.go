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

var id int       // Temporary storage of id value of table journal_entries
var date string  // Temporary storage of date value of table journal_entries
var entry string // Temporary storage of entry value of table journal_entries

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

	rows, err := db.Query(`SELECT * FROM journal_entries WHERE date = ?`, journalDate)
	if err != nil {
		log.Fatal(err)
	}

	dateExists := false

	for rows.Next() {
		err := rows.Scan(&id, &date, &entry)
		if err != nil {
			log.Fatal(err)
		}
		if journalDate == date {
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
		rows.Scan(&id, &date, &entry)

		journalEntry = fmt.Sprint(entry + "\n\n" + journalEntry)

		statement, err := db.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec(journalEntry, journalDate)

	} else {
		statement, err := db.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec(journalDate, journalEntry)

	}

	rows, err = db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&id, &date, &entry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(date + ":\n" + entry)
	}
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

	rows, err := db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	if err != nil {
		log.Fatal(err)
	}

	dateExists := false

	for rows.Next() {
		err := rows.Scan(&id, &date, &entry)
		if err != nil {
			log.Fatal(err)
		}
		if journalDate == date {
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
		rows.Next()
		rows.Scan(&id, &date, &entry)
		journalEntry = fmt.Sprint(entry + "\n\n" + journalEntry)

		statement, err := db.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec(journalEntry, journalDate)

	} else {
		statement, err := db.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec(journalDate, journalEntry)
	}

	rows, err = db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&id, &date, &entry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(date + ":\n" + entry)
	}
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

	for rows.Next() {
		err := rows.Scan(&id, &date, &entry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(date + ":\n" + entry)
	}
}

// ViewEntireJournal prints every date and entry of journal_entries
func ViewEntireJournal(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM journal_entries ORDER BY date")

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&id, &date, &entry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(date + ":\n" + entry)
	}
}

// DeleteEntry deletes the record of a particular date
func DeleteEntry(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to delete (MM-DD-YYYY):")
	journalDate, err := reader.ReadString('\n')
	journalDate = journalDate[:len(journalDate)-1]

	if err != nil {
		log.Fatal(err)
	}

	statement, err := db.Prepare("DELETE FROM journal_entries WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
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
	for rows.Next() {
		err := rows.Scan(&id, &date, &entry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(date + ":\n" + entry)
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
	statement.Exec(journalEntry, journalDate)

	rows, err = db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&id, &date, &entry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(date + ":\n" + entry)
	}
}
