package main

import (
        "fmt"
        "golobsters/lobsterdb"
        "code.google.com/p/gosqlite/sqlite"
        "log"
        "os"
)

func read_stories(filename string) []string {
        db, err := sqlite.Open(filename)
        if err != nil {
                log.Fatal("[!] couldn't open filename")
        }

        defer db.Close()

        res, err := db.Exec("select guid from posted")
        if err != nil {
                log.Fatal("[!] could select from posted")
        }

        n, _ := res.RowsAffected()
        guids := make([]string, n) 

        return guids
}

func mark_posted(guids []string) bool {
        errs := 0
        for _, guid := range guids {
                if err := lobsterdb.PostStory(guid); err != nil {
                        errs++
                        log.Println("[!] error posting story ", guid)
                } else {
                        fmt.Printff("[+] marking %s as posted.", guid)
                }
        }

        return errs == 0
}

func main() {
        if len(os.Args) == 1 {
                log.Fatal("no filename specified")
        }

        fmt.Println(read_stories(os.Args[1]))
}
