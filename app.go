/*
   golobsters is an application that checks for new posts on lobste.rs
   and posts new posts to the corresponding Twitter account.
*/

package main

import (
	"fmt"
	"github.com/gokyle/golobsters/bot"
	"github.com/gokyle/golobsters/frontend"
	"github.com/gokyle/gomon/monitor"
	"log"
	"net/http"
	"os"
	"time"
)

var siteURL = "http://lobsternews.kyleisom.net"

// run is stubbed
func validate_env_var(name string) bool {
	value := os.Getenv(name)
	if value == "" {
		return false
	}
	return true
}

// should check to ensure required environment variables are present
func validate_environment() {
	vars := []string{"TW_CKEY", "TW_CSEC", "TW_ATOK", "TW_ASEC",
		"PG_DBNAME", "PG_USER", "PG_PASS", "PG_HOST",
		"PG_PORT", "PG_SSLMODE"}
	for _, name := range vars {
		if !validate_env_var(name) {
			panic(fmt.Sprintf("missing environment variable: %s",
				name))
		}
	}
}

func herokuPing() {
	for {
		log.Println("pinging site")
		http.Get(siteURL)
		<-time.After(15 * time.Minute)
	}
}

func init() {
	log.Println("[+] initialising application")
	validate_environment()
	monitor.ConfigFromEnv()
	if !monitor.EmailEnabled() {
		log.Fatal("[!] error configuring monitor: mail not configured")
	}

	if !monitor.PushoverEnabled() {
		log.Fatal("[!] error configuring monitor: pushover not configured")
	}
}

func main() {
	log.Println("[+] app starting")
	go monitor.Monitor(bot.Run)
	log.Println("[+] launching web front end")
	frontend.HttpServer()
	log.Println("[+] app stopping")
}
