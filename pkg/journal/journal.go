// Package journal connects the user to the journal database and allows
// the user to alter the table journal_entries.
package journal

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

var dbid int       // id value of table journal_entries
var dbdate string  // date value of table journal_entries
var dbentry string // entry value of table journal_entries

// InputEntry adds the current date as a string and prompts the user for
// a journal entry input to be stored into the database in association
// with the date.
func InputEntry(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nAdd An Entry For Today\n\n")

	journalDate := string(time.Now().Format("01-02-2006"))

	fmt.Printf("Input journal entry for %s:\n", journalDate)
	journalEntry, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	journalEntry = journalEntry[:len(journalEntry)-1]

	ifEntryExists(db, journalEntry, journalDate)

	printEntry(db, journalDate)
}

// InputEntryDate prompts the user for a date as a string and prompts the user
// for a journal entry input to be stored into the database in association
// with the date.
func InputEntryDate(db *sql.DB) {
	var journalDate string
	var err error

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nAdd An Entry For Another Date\n\n")

	for {
		fmt.Println("Input date (MM-DD-YYYY): ")
		journalDate, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		journalDate = journalDate[:len(journalDate)-1]

		// Checks to see if the inputted date is in the correct format
		matched, err := regexp.MatchString(`(0[1-9]|1[012])[- /.](0[1-9]|[12][0-9]|3[01])[- /.](19|20)[0-9][0-9]`, journalDate)
		if err != nil {
			log.Fatal(err)
		}
		if matched == true {
			break
		}
		fmt.Println("Incorrect date format. Please try again.")
	}

	fmt.Printf("Input journal entry for %s:\n", journalDate)
	journalEntry, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	journalEntry = journalEntry[:len(journalEntry)-1]

	ifEntryExists(db, journalEntry, journalDate)

	printEntry(db, journalDate)
}

// ViewEntry prints the date and entry of a particular date
func ViewEntry(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nView An Entry\n\n")

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
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&dbid, &dbdate, &dbentry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dbdate + ":\n" + dbentry)
	}
}

// ViewEntireJournal prints every date and entry of journal_entries
func ViewEntireJournal(db *sql.DB) {
	fmt.Print("\nView All Entries\n\n")

	rows, err := db.Query("SELECT * FROM journal_entries ORDER BY date")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&dbid, &dbdate, &dbentry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\n" + dbdate + ":\n" + dbentry)
	}
}

// DeleteEntry deletes the record of a particular date
func DeleteEntry(db *sql.DB) {
	var journalDate string
	var err error

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nDelete An Entry\n")

	for {
		fmt.Println("Input date of journal entry to delete (MM-DD-YYYY):")
		journalDate, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		journalDate = journalDate[:len(journalDate)-1]

		// Checks to see if the inputted date is in the correct format
		matched, err := regexp.MatchString(`(0[1-9]|1[012])[- /.](0[1-9]|[12][0-9]|3[01])[- /.](19|20)[0-9][0-9]`, journalDate)
		if err != nil {
			log.Fatal(err)
		}
		if matched == true {
			break
		}
		fmt.Println("Incorrect date format. Please try again.")
	}

	statement, err := db.Prepare("DELETE FROM journal_entries WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	statement.Exec(journalDate)

	fmt.Printf("\nEntry for %s has been deleted.\n", journalDate)
}

// DeleteJournal deletes the entire table of journal_entries
func DeleteJournal(db *sql.DB, username string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nDelete Entire Journal\n\n")

	for {
		fmt.Print("Are you sure you want to delete your entire journal? (Y/n): ")

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
				fmt.Println("\nExiting GoJournal")
				os.Exit(0)
			}
		}
		fmt.Println("Not a valid choice. Please try again.")
	}

	dataSource := fmt.Sprintf("./databases/%s.db", username)
	err := os.Remove(dataSource)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nJournal for %s has been deleted.\n", username)
}

// EditEntry replaces the entry of a particular date
func EditEntry(db *sql.DB) {
	var journalDate string
	var err error

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nEdit an Entry\n\n")

	for {
		fmt.Println("Input date of journal entry to edit:")
		journalDate, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		journalDate = journalDate[:len(journalDate)-1]

		// Checks to see if the inputted date is in the correct format
		matched, err := regexp.MatchString(`(0[1-9]|1[012])[- /.](0[1-9]|[12][0-9]|3[01])[- /.](19|20)[0-9][0-9]`, journalDate)
		if err != nil {
			log.Fatal(err)
		}
		if matched == true {
			break
		}
		fmt.Println("Incorrect date format. Please try again.")
	}

	printEntry(db, journalDate)

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
	defer statement.Close()
	statement.Exec(journalEntry, journalDate)

	printEntry(db, journalDate)
}

// ifEntryExists checks to see if an entry for a certain date already exists
// in the table journal_entires in a specified SQL database.
func ifEntryExists(db *sql.DB, journalEntry string, journalDate string) {
	rows, err := db.Query(`SELECT * FROM journal_entries WHERE date = ?`, journalDate)
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
		rows, err = db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
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

		statement, err := db.Prepare("UPDATE journal_entries SET entry = ? WHERE date = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()
		statement.Exec(journalEntry, journalDate)

	} else {
		statement, err := db.Prepare("INSERT INTO journal_entries (date, entry) VALUES (?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()
		statement.Exec(journalDate, journalEntry)
	}
}

// printEntry prints the entry of a specified date onto the console
func printEntry(db *sql.DB, journalDate string) {
	rows, err := db.Query("SELECT * FROM journal_entries WHERE date = ?", journalDate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&dbid, &dbdate, &dbentry)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\n" + dbdate + ":\n" + dbentry)
	}
}
