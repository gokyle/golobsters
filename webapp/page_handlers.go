package webapp

import (
        "fmt"
        "golobsters/bot"
        "net/http"
)

func rootPage(w http.ResponseWriter, req *http.Request) {
        page := "last update: " + bot.LastUpdate()
        fmt.Fprintln(w, page) 
}
