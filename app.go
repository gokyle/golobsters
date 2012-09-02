/*
   golobsters is an application that checks for new posts on lobste.rs
   and posts new posts to the corresponding Twitter account.
*/

package main

import (
	"fmt"
	"github.com/gokyle/gomon/monitor"
	"golobsters/bot"
        "golobsters/webapp"
	"log"
	"os"
	"time"
)

// run is stubbed
func run() {
	go monitor.Monitor(bot.Run)
	time.Sleep(5 * 1000 * time.Millisecond)
	log.Println("[+] bot last update: ", bot.LastUpdate)
        for {
       	        if "" != bot.LastUpdate() {
		        log.Printf("[+] bot last update: %s\n", bot.LastUpdate())
	        }

                time.Sleep(15 * 1000 * time.Millisecond)
        }
	return
}

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

func init() {
	validate_environment()
	err := monitor.ConfigFromJson()
	if err != nil {
		fmt.Println("[!] error configuring monitor: ", err)
		os.Exit(1)
	}
}

func main() {
	go run()
        webapp.HttpServer()
}


