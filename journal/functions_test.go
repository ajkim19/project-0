package journal

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestViewEntireJournal(t *testing.T) {
	database, err := sql.Open("sqlite3", "./journal_test.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, err := database.Prepare("DROP TABLE IF EXISTS journal_entries")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-01-2020", "Today is New Year's Day!")`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-02-2020", "It's the day after New Year's Day!")`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ViewEntireJournal(database)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != "01-01-2020: Today is New Year's Day!\n01-02-2020: It's the day after New Year's Day!\n" {
		t.Error("Expected 01-01-2020: Today is New Year's Day!\n01-02-2020: It's the day after New Year's Day!\n")
	}
}

// func TestDeleteJournal(t *testing.T) {
// 	database, err := sql.Open("sqlite3", "./journal_test.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	statement, err := database.Prepare("DROP TABLE IF EXISTS journal_entries")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	statement.Exec()

// 	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	statement.Exec()

// 	statement, err = database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-01-2020", "Today is New Year's Day!")`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	statement.Exec()

// 	statement, err = database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-02-2020", "It's the day after New Year's Day!")`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	statement.Exec()

// 	rescueStdout := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w

// 	DeleteJournal(database)

// 	w.Close()
// 	out, _ := ioutil.ReadAll(r)
// 	os.Stdout = rescueStdout

// 	if string(out) != "" {
// 		t.Error("Expected no output")
// 	}
// }
