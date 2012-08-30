package main

import (
        "fmt"
        "golobsters/lobsterdb"
        _ "github.com/mattn/go-sqlite3"
	"database/sql"
        "log"
        "os"
)

func read_stories(filename string) []string {
        db, err := sql.Open("sqlite3", filename)
        if err != nil {
                log.Fatal("[!] couldn't open filename: ", err)
        }

        defer db.Close()

        rows, err := db.Query("select guid from posted")
        if err != nil {
                log.Fatal("[!] could select from posted")
        }

        guids := make([]string, 0) 

        for rows.Next() {
                var guid string
                rows.Scan(&guid)
                guids = append(guids, guid)
        }

        return guids
}

func mark_posted(guids []string) bool {
        errs := 0
        fmt.Printf("[+] attempting to mark %d stories as posted...", len(guids))
        for _, guid := range guids {
                if err := lobsterdb.PostStory(guid); err != nil {
                        errs++
                        log.Println("[!] error posting story ", guid)
                } else {
                        fmt.Printf("[+] marking %s as posted.", guid)
                }
        }

        return errs == 0
}

func main() {
        if len(os.Args) == 1 {
                fmt.Println("[!] no filename specified")
                os.Exit(1)
        }

        guids := read_stories(os.Args[1])
        if len(guids) == 0 {
                log.Fatal("could not retrieve from the database")
        }

        if !mark_posted(guids) {
                log.Printf("[!] error importing into postgres")
        }
}
