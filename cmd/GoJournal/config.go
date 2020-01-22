package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	username string  // Username of journal
	date     bool    // Flag to input an entry for a particular date
	view     bool    // Flag to view table journal_entries
	delete   bool    // Flag to delete from table journal_entries
	edit     bool    // Flag to edit journal_entries
	all      bool    // Flag to apply alteration to the entire table of journal_entries
	flag1    string  // First flag string
	flag2    string  // Second flag string
	database *sql.DB // Pointer to database handle
	err      error   // Temporary storage of error value
	dbid     int     // Temporary storage of id value of table journal_entries
	dbdate   string  // Temporary storage of date value of table journal_entries
	dbentry  string  // Temporary storage of entry value of table journal_entries
)

func init() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(" _______________________")
	fmt.Println("|                       |")
	fmt.Println("| Welcome to GoJournal! |")
	fmt.Print("|_______________________|\n\n")

	// Prompts for username
	for {
		fmt.Print("Please enter username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		username = username[:len(username)-1]

		if strings.ContainsAny(username, " !@#$%^&*()[]{}`~:;<>,./\\+*\"?'") == false {
			break
		}

		fmt.Println("Invalid username. Please Try again")
	}

	// Makes a handle for the database journal
	dataSourceName := fmt.Sprintf("./databases/%s.db", username)
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

	rows, err := database.Query("SELECT * FROM journal_entries")
	if err != nil {
		log.Fatal(err)
	}

	var dateExists bool

	for rows.Next() {
		rows.Scan(&dbid, &dbdate, &dbentry)
		if dbdate == "01-06-2020" {
			dateExists = true
		}
	}

	// Adds an entry to journal_entries if it is empty
	if dateExists == false {
		statement, err := database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-06-2020", "Welcome to Revature!")`)
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec()
	}

	// Initalizes the flags
	flag.BoolVar(&date, "date", false, "add entry to specified date")
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
