package main

import (
	"fmt"
	"net/http"

	"github.com/ajkim19/project-0/GoJournal/pkg/templates"
)

func main() {
	println("Server is running on port 8080")

	http.Handle("/", http.FileServer(http.Dir("html")))
	http.HandleFunc("/GoJournal", func(w http.ResponseWriter, r *http.Request) {
		var entry string = r.FormValue("entry")
		if entry == "" {
			entry = "Random Entry"
		}
		fmt.Fprint(w, templates.OpeningHTML)
		fmt.Fprint(w, templates.Head)
		fmt.Fprint(w, templates.OpeningBody)
		fmt.Fprint(w, "<h1>", entry, "</h1>")
		fmt.Fprint(w, templates.ClosingBody)
		fmt.Fprint(w, templates.ClosingHTML)
	})
	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {})
	http.ListenAndServe(":8080", nil)
}
