package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	journalEntries := map[string]string{}

	reader := bufio.NewReader(os.Stdin)

	journalDate := string(time.Now().Format("01-02-2006"))

	fmt.Println("Input Journal Entry:")
	text, _ := reader.ReadString('\n')

	journalEntries[journalDate] = text

	fmt.Println(journalEntries)
}
