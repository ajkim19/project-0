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

	"github.com/ajkim19/project-0/pkg/journal"
	_ "github.com/mattn/go-sqlite3"
)

var (
	reg      *regexp.Regexp // Pointer to Regexp object
	username string         // Username of journal
	help     bool           // Flag to print help menu
	date     bool           // Flag to input an entry for a particular date
	view     bool           // Flag to view table journal_entries
	delete   bool           // Flag to delete from table journal_entries
	edit     bool           // Flag to edit journal_entries
	all      bool           // Flag to apply alteration to the entire table of journal_entries
	flag1    string         // First flag string
	flag2    string         // Second flag string
	database *sql.DB        // Pointer to database handle
	err      error          // Temporary reference to error value
	dbid     int            // Temporary reference to id value of table journal_entries
	dbdate   string         // Temporary reference to date value of table journal_entries
	dbentry  string         // Temporary reference to entry value of table journal_entries
)

func init() {
	// Initalizes the flags
	flag.BoolVar(&help, "help", false, "prints help menu")
	flag.BoolVar(&date, "date", false, "add entry to specified date")
	flag.BoolVar(&view, "view", false, "view entry")
	flag.BoolVar(&delete, "delete", false, "delete entry")
	flag.BoolVar(&edit, "edit", false, "edit entry")
	flag.BoolVar(&all, "all", false, "apply to every entry")
	flag.Parse()

	// Removes special characters of flags
	reg, err = regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 2 {
		flag1 = reg.ReplaceAllString(os.Args[1], "")
	} else if len(os.Args) > 2 {
		flag1 = reg.ReplaceAllString(os.Args[1], "")
		flag2 = reg.ReplaceAllString(os.Args[2], "")
	}

	// Initialize reader
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(" _______________________")
	fmt.Println("|                       |")
	fmt.Println("| Welcome to GoJournal! |")
	fmt.Print("|_______________________|\n\n")

	// Checks for help flag
	if flag1 == "help" {
		journal.Help()
	}

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

	dataSourceName := fmt.Sprintf("./databases/%s.db", username)

	// Detects if database file exists
	var _, err = os.Stat(dataSourceName)

	if os.IsNotExist(err) {
		fmt.Println("\nA journal for this username does not exist.")
		for {
			fmt.Print("Would you like to create one (Y/n): ")
			choice, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			choice = choice[:len(choice)-1]

			// Checks to see if the input is valid
			matched, err := regexp.MatchString(`[Y]|[n]`, choice)
			if err != nil {
				log.Fatal(err)
			}
			if matched == true {
				if choice == "Y" {
					break
				} else {
					os.Exit(0)
				}
			}
			fmt.Println("Not a valid choice. Please try again.")
		}
	} else if err != nil {
		log.Fatal(err)
	}

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

	rows, err := database.Query("SELECT * FROM journal_entries")
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
}
