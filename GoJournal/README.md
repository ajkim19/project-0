# GoJournal
## Author: Aaron Kim

## About
GoJournal is a journal program written in Go for inputing, editing, viewing, and deleting journal entries. It uses SQLite for storage of dates and entries..

## How to Use GoJournal
When using GoJournal, flags may be used to utilize certain functions.
Flags:
- default (no flags) - Allows you to input a journal entry to the date it is written. 
- "date" - Allows you to specify a date to your journal entry
- "edit" - Allows you to edit an existing journal entry at a specified date.
- "view" - Allows you to view an existing journal entry at a specified date.
- "delete" - Allows you to delete an existing journal entry at a specified date.
- "all" - When following a "view" or a "delete" flag, the followed feature will apply to the entire journal.
