package frontend

import (
	"fmt"
	"golobsters/bot"
	"golobsters/lobsterdb"
	"net/http"
)

func rootPage(w http.ResponseWriter, req *http.Request) {
        page := "twitter account: @lobsternews\n"
        page += "git repo: git clone git://github.com/gokyle/golobsters.git\n"
        page += "github page: http://gokyle.github.com/golobsters/\n\n"
        page += "stats\n=====\n"
	page += "last update: " + bot.LastUpdate()
	page += fmt.Sprintf("\nstories posted: %d\n", lobsterdb.CountStories())
	fmt.Fprintln(w, page)
}
