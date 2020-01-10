package main

import (
	"strconv"
	"bufio"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"os"
	"database/sql"
	"time"
)

func main() {
	database, _ := sql.Open(driverName: "sqlite3", dataSourceName: "./journal.db")
	statement, _ := database.Prepare(query: "CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	statement.Exec()
	statement, _ = database.Prepare(query: "INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
	statement.Exec(args: "01-01-2020", "It is New Year's Day today.")
	rows, _ := database.Query(query: "SELECT * FROM journal_entries")
	var id int
	var date string
	var entry string
	for rows.Next() {
		rows.Scan(&id, &date, &entry)
		fmt.Println(strconv.Itoa(id) + ": " + date + " " + entry)
	}
	

	//Creates map for journal entries
	journalEntries := map[string]string{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Input Journal Entry:")
	text, _ := reader.ReadString('\n')
	journalDate := string(time.Now().Format("01-02-2006"))
	journalEntries[journalDate] = text

	fmt.Println(journalEntries)
}
