package main

import (
        "fmt"
        "golobsters/lobsterdb"
        "googlecode.com/hg/sqlite"
        "log"
)

func read_stories(filename string) []string {
        db, err = sqlite.Open(filename)
        if err != nil {
                log.Fatal("[!] couldn't open filename")
        }

        defer db.Close()

        res, err := sqlite.Exec("select guid from posted")
        if err != nil {
                log.Fatal("[!] could select from posted")
        }

        guids := (make []string, res.RowsAffected())

        return guids
}

func main() {
        if len(os.Args) == 1 {
                log.Fatal("no filename specified")
        }

        fmt.Println(read_stories(os.Args[1])
}
