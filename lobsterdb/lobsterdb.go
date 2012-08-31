// lobsterdb implements the database interactivity for the lobster bot.
package lobsterdb

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"log"
	"os"
)

// ConnStringFromEnv loads the database credentials from the environment. 
func ConnStringFromEnv() string {
	return fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%s sslmode=%s",
		os.Getenv("PG_DBNAME"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_SSLMODE"))
}

// StoryPosted is used to determine whether a story has been posted or not. It 
// is keyed to the story's id url; for example, "https://lobste.rs/s/lwrxft/"
func StoryPosted(guid string) (bool, error) {
	db, err := sql.Open("postgres", ConnStringFromEnv())
	if err != nil {
		log.Printf("[!] lobsterdb couldn't open database connection: %s",
			err)
		return true, err
	}
	defer db.Close()
        log.Println("[+] lobsterdb connected to database (preparing select)")

	res, err := db.Exec("select posted from posted where guid=$1")
	if err != nil {
		log.Printf("[!] lobsterdb select error: %s", err)
		return true, err
	}

	if n, _ := res.RowsAffected(); n > 1 {
		log.Printf("[!] lobsterdb %s has more than one row", guid)
	} else if n, _ := res.RowsAffected(); n == 0 {
		return false, nil
	}

	return true, err
}

// PostStory is used to mark a story as posted in the database.
func PostStory(guid string) error {
	db, err := sql.Open("postgres", ConnStringFromEnv())
	if err != nil {
		log.Printf("[!] lobsterdb couldn't open database connection: %s",
			err)
		return err
	}
	defer db.Close()
        log.Printf("[+] lobsterdb connected to database (preparing insert)")


        res, err := db.Exec("insert into posted (guid, posted) values ($1, $2)",
                            guid, true)
        if err != nil {
                log.Printf("[!] lobsterdb couldn't insert into database",
                    guid)
                return err
        }

        if n, _ := res.RowsAffected(); n == 0 {
                log.Printf("[!] lobsterdb insert affects 0 rows")
                return fmt.Errorf("insert affects 0 rows")
        }

        return nil
}
