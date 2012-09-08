package frontend

import (
	"fmt"
	"github.com/gokyle/golobsters/bot"
	"github.com/gokyle/golobsters/dbase"
	"net/http"
)

func rootPage(w http.ResponseWriter, req *http.Request) {
	db, err := dbase.ConnectFromEnv()
	stats := ""

	if err == nil {
		stats += "stats\n=====\n"
		stats += "last tweet: " + bot.LastUpdate()
		stats += fmt.Sprintf("\nstories posted: %d\n", dbase.CountStories(db))
	} else {
		stats += "couldn't connect to database: " + err.Error()
	}

	page := "twitter account: @lobsternews\n"
	page += "git repo: git clone git://github.com/gokyle/golobsters.git\n"
	page += "github page: http://gokyle.github.com/golobsters/\n\n"
	page += stats
	fmt.Fprintln(w, page)
}