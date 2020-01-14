package journal

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"time"
)

var id int
var date string
var entry string

// InputEntry adds journal entry into database
func InputEntry(d *sql.DB) {
	// Prompts user for journal input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input journal entry:")
	journalEntry, _ := reader.ReadString('\n')
	journalEntry = journalEntry[:len(journalEntry)-1]

	// Grabs the date of the entry made
	journalDate := string(time.Now().Format("01-02-2006"))

	rows, _ := d.Query("SELECT * FROM journal_entries")
	dateExists := false

	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		if journalDate == date {
			dateExists = true
		}
	}

	if dateExists {
		// Adds entry onto the entry with the same date into the database if an entry for the date already exists
		rows, _ = d.Query("SELECT * FROM journal_entries")
		rows.Scan(&id, &date, &entry)
		journalEntry = fmt.Sprint(entry + "\n\n" + journalEntry)

		statement, _ := d.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
		statement.Exec(journalEntry, journalDate)

	} else {
		// Inserts date and entry into the database if an entry for the date does not already exists
		statement, _ := d.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
		statement.Exec(journalDate, journalEntry)
	}

	rows, _ = d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	rows.Scan(&id, &date, &entry)
	fmt.Println(date + ": " + entry)
}

// ViewEntry prints the entry of a particular date
func ViewEntry(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to view:")
	journalDate, _ := reader.ReadString('\n')
	journalDate = journalDate[:len(journalDate)-1]

	rows, _ := d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)

	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(date + ": " + entry)
	}
}

// ViewEntireJournal prints the entire table of journal_entries
func ViewEntireJournal(d *sql.DB) {
	rows, _ := d.Query("SELECT * FROM journal_entries")
	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(date + ": " + entry)
	}
}

// DeleteEntry deletes the entry of a particular date
func DeleteEntry(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to delete:")
	journalDate, _ := reader.ReadString('\n')

	statement, _ := d.Prepare("DELETE FROM journal_entries WHERE date = ?")
	statement.Exec(journalDate)
}

// DeleteTable deletes the entire table of journal_entries
func DeleteTable(d *sql.DB) {
	statement, _ := d.Prepare("DROP TABLE journal_entries")
	statement.Exec()
}

// EditEntry replaces the entry of a particular date
func EditEntry(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to edit:")
	journalDate, _ := reader.ReadString('\n')
	journalDate = journalDate[:len(journalDate)-1]

	rows, _ := d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(date + ": " + entry)
	}

	fmt.Println("Input replacement entry:")
	journalEntry, _ := reader.ReadString('\n')
	journalEntry = journalEntry[:len(journalEntry)-1]

	statement, _ := d.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
	statement.Exec(journalEntry, journalDate)

	rows, _ = d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(date + ": " + entry)
	}
}
