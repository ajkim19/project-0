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

	// Prompts user for journal input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input Journal Entry:")
	journalEntry, _ := reader.ReadString('\n')

	// Grabs the date of the entry made
	journalDate := string(time.Now().Format("01-02-2006"))

	// Inserts date and entry into the database
	statement, _ = database.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
	statement.Exec(journalDate, journalEntry)

	// printEntireJournal(database)
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
