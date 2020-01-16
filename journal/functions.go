package journal

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"time"
	"log"
)

var id int
var date string
var entry string

// InputEntry adds the current date as a string and prompts user for
// a journal entry input to be stored into a sql database in association
// with the date. If an entry for the specified date already exists,
// the preexisting entry will be printed out and the user will be prompted
// to add to the entry. A new line and the entry will then be concatenated
// to the preexisting entry. Lastly, InputEntry prints the date and
// the complete entry.
func InputEntry(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input journal entry:")
	journalEntry, err1 := reader.ReadString('\n')\
	if err1 != nil {
		log.Fatal(err1)
	}
	journalEntry = journalEntry[:len(journalEntry)-1]

	journalDate := string(time.Now().Format("01-02-2006"))

	rows, err2 := d.Query("SELECT * FROM journal_entries")
	if err2 != nil {
		log.Fatal(err2)
	}
	dateExists := false

	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		if journalDate == date {
			dateExists = true
		}
	}

	if dateExists {
		rows, err2 = d.Query("SELECT * FROM journal_entries")
		if err2 != nil {
			log.Fatal(err2)
		}
		rows.Scan(&id, &date, &entry)
		journalEntry = fmt.Sprint(entry + "\n\n" + journalEntry)

		statement, err3 := d.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
		if err3 != nil {
			log.Fatal(err3)
		}
		statement.Exec(journalEntry, journalDate)

	} else {
		statement, err3 := d.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
		if err3 != nil {
			log.Fatal(err3)
		}
		statement.Exec(journalDate, journalEntry)
	}

	rows, _ = d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	rows.Next()
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

	rows.Next()
	rows.Scan(&id, &date, &entry)
	fmt.Println(date + ": " + entry)
}

// ViewEntireJournal prints the entire table of journal_entries
func ViewEntireJournal(d *sql.DB) {
	rows, err := d.Query("SELECT * FROM journal_entries")
	if err != nil {
		log.Fatal(err)
	}

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

	statement, err := d.Prepare("DELETE FROM journal_entries WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(journalDate)
	
}

// DeleteTable deletes the entire table of journal_entries
func DeleteTable(d *sql.DB) {
	statement, err := d.Prepare("DROP TABLE journal_entries")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

// EditEntry replaces the entry of a particular date
func EditEntry(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to edit:")
	journalDate, _ := reader.ReadString('\n')
	journalDate = journalDate[:len(journalDate)-1]

	rows, _ := d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	rows.Next()
	rows.Scan(&id, &date, &entry)
	fmt.Println(date + ": " + entry)

	fmt.Println("Input replacement entry:")
	journalEntry, _ := reader.ReadString('\n')
	journalEntry = journalEntry[:len(journalEntry)-1]

	statement, _ := d.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
	statement.Exec(journalEntry, journalDate)

	rows, _ = d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	rows.Next()
	rows.Scan(&id, &date, &entry)
	fmt.Println(date + ": " + entry)
}
