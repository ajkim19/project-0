package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Opens database connection
	database, _ := sql.Open("sqlite3", "./journal.db")

	// Creates table if it does not exist
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	statement.Exec()

	//inputEntry(database)
	//printEntireJournal(database)
	//searchEntry(database)
	editEntry(database)
}

// Adds journal entry into database
func inputEntry(d *sql.DB) {
	// Prompts user for journal input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input journal entry:")
	journalEntry, _ := reader.ReadString('\n')
	journalEntry = journalEntry[:len(journalEntry)-1]

	// Grabs the date of the entry made
	journalDate := string(time.Now().Format("01-02-2006"))

	// Inserts date and entry into the database
	statement, _ := d.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
	statement.Exec(journalDate, journalEntry)
}

// Prints entire table of journal_entries
func printEntireJournal(d *sql.DB) {
	rows, _ := d.Query("SELECT * FROM journal_entries")
	var id int
	var date string
	var entry string
	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(strconv.Itoa(id) + ": " + date + " " + entry)
	}
}

func deleteTable(d *sql.DB) {
	statement, _ := d.Prepare("DROP TABLE journal_entries")
	statement.Exec()
}

func deleteEntry(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to delete:")
	journalDate, _ := reader.ReadString('\n')

	statement, _ := d.Prepare("DELETE FROM journal_entries WHERE date = ?")
	statement.Exec(journalDate)
}

func searchEntry(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to view:")
	journalDate, _ := reader.ReadString('\n')
	journalDate = journalDate[:len(journalDate)-1]

	rows, _ := d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)

	var id int
	var date string
	var entry string
	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(strconv.Itoa(id) + ": " + date + " " + entry)
	}
}

func editEntry(d *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input date of journal entry to edit:")
	journalDate, _ := reader.ReadString('\n')
	journalDate = journalDate[:len(journalDate)-1]

	rows, _ := d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	var id int
	var date string
	var entry string
	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(strconv.Itoa(id) + ": " + date + " " + entry)
	}

	fmt.Println("Input replacement entry:")
	journalEntry, _ := reader.ReadString('\n')
	journalEntry = journalEntry[:len(journalEntry)-1]
	fmt.Println(journalEntry)

	statement, _ := d.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
	statement.Exec(journalEntry, journalDate)

	rows, _ = d.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(strconv.Itoa(id) + ": " + date + " " + entry)
	}
}
