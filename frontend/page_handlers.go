package frontend

import (
	"fmt"
	"golobsters/bot"
	"golobsters/lobsterdb"
	"net/http"
)

func rootPage(w http.ResponseWriter, req *http.Request) {
	page := "last update: " + bot.LastUpdate()
	page += fmt.Sprintf("\nstories posted: %d\n", lobsterdb.CountStories)
	fmt.Fprintln(w, page)
}
