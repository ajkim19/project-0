package journalhtml

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ajkim19/project-0/pkg/templates"
	"github.com/ajkim19/project-1/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
)

type jEntry struct {
	Date  string
	Entry string
}

var jEntries = []jEntry{}

var database *sql.DB // Pointer to database handle
var err error        // Temporary reference to error value
var dbid int         // Temporary reference to id value of table journal_entries
var dbdate string    // Temporary reference to date value of table journal_entries
var dbentry string   // Temporary reference to entry value of table journal_entries
var rows *sql.Rows
var username string = "journal"
var revproxyauth = os.Getenv("REVPROXYAUTH")
var basicAuth string

func init() {
	basicAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte(revproxyauth))

	dataSourceName := fmt.Sprintf("./GoJournal/databases/%s.db", username)

	// Makes a handle for the database journal
	database, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Creates the table journal_entries if it has been dropped
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	rows, err = database.Query("SELECT * FROM journal_entries")
	if err != nil {
		log.Fatal(err)
	}

	var dateExists bool

	for rows.Next() {
		rows.Scan(&dbid, &dbdate, &dbentry)
		if dbdate == "2020-01-06" {
			dateExists = true
		}
	}

	// Adds an entry to journal_entries if it is empty
	if dateExists == false {
		statement, err := database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("2020-01-06", "Welcome to Revature!")`)
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec()
	}

	rows, err = database.Query("SELECT * FROM journal_entries ORDER BY date DESC")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&dbid, &dbdate, &dbentry)
		jEntries = append(jEntries, jEntry{Date: dbdate, Entry: dbentry})
	}
}

func GoJournalHTML(w http.ResponseWriter, r *http.Request) {

	requestAddr := r.RemoteAddr

	if r.Header.Get("Proxy-Authorization") != basicAuth {
		logger.Logger.Printf("Unauthorized client: %s\n", requestAddr)
	}

	w.Header().Set("Content-Type", "text/html")

	var journalDate string = r.FormValue("date")

	// Checks to see if the inputted date is in the correct format
	matched, err := regexp.MatchString(`(0[1-9]|1[012])[- /.](0[1-9]|[12][0-9]|3[01])[- /.](19|20)[0-9][0-9]`, journalDate)
	if err != nil {
		log.Fatal(err)
	}
	if matched == false {
		journalDate = string(time.Now().Format("2006-01-02"))
	}

	var journalEntry string = r.FormValue("entry")
	if journalEntry != "" {
		rows, err = database.Query(`SELECT * FROM journal_entries WHERE date = ?`, journalDate)
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
			rows, err = database.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
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

			statement, err := database.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
			if err != nil {
				log.Fatal(err)
			}
			defer statement.Close()
			statement.Exec(journalEntry, journalDate)

		} else {
			statement, err := database.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			defer statement.Close()
			statement.Exec(journalDate, journalEntry)
		}
	}

	rows, err = database.Query("SELECT * FROM journal_entries ORDER BY date DESC")
	if err != nil {
		log.Fatal(err)
	}

	jEntries = []jEntry{}

	for rows.Next() {
		rows.Scan(&dbid, &dbdate, &dbentry)
		jEntries = append(jEntries, jEntry{Date: dbdate, Entry: dbentry})
	}

	fmt.Fprint(w, templates.OpeningHTML)
	fmt.Fprint(w, templates.Head)
	fmt.Fprint(w, templates.OpeningBody)
	for _, v := range jEntries {
		fmt.Fprint(w, "<h2>", v.Date, "</h2>")
		for _, line := range strings.Split(strings.TrimSuffix(v.Entry, "\n"), "\n") {
			fmt.Fprint(w, "<h2>", line, "</h2>")
		}
	}
	fmt.Fprint(w, templates.ClosingBody)
	fmt.Fprint(w, templates.ClosingHTML)

}
