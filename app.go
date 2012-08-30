/*
   golobsters is an application that checks for new posts on lobste.rs
   and posts new posts to the corresponding Twitter account.
*/

package main

import (
	"fmt"
	"golobsters/lobsterdb"
	"os"
)

// run is stubbed
func run() {
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
		"PG_PORT", "PG_SSL"}
	for _, name := range vars {
		if !validate_env_var(name) {
			panic(fmt.Sprintf("missing environment variable: %s",
				name))
		}
	}
}

func main() {
	fmt.Println(lobsterdb.ConnStringFromEnv())
	run()
}
