// dbase implements the database interactivity for the lobster bot.
package dbase

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"log"
	"os"
)

// ConnStringFromEnv loads the database credentials from the environment.
func ConnectFromEnv() (*sql.DB, error) {
	conn_string := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%s sslmode=%s",
		os.Getenv("PG_DBNAME"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_SSLMODE"))
	db, err := sql.Open("postgres", conn_string)
	if err != nil {
		log.Printf("[!] dbase couldn't open database connection: %s",
			err)
		db = nil
	}

	return db, err
}

// StoryPosted is used to determine whether a story has been posted or not. It
// is keyed to the story's id url; for example, "https://lobste.rs/s/lwrxft/"
func StoryPosted(db *sql.DB, guid string) (bool, error) {

	rows, err := db.Query("select posted from posted where guid=$1", guid)
	if err != nil {
		log.Printf("[!] dbase select error: %s", err)
		return true, err
	}

	log.Println("[+] dbase select query completed, retrieving results")
	row_count := 0
	for rows.Next() {
		row_count++
	}

	if row_count > 1 {
		log.Printf("[!] dbase %s has more than one row", guid)
	} else if row_count == 0 {
		return false, nil
	}

	return true, err
}

// PostStory is used to mark a story as posted in the database.
func PostStory(db *sql.DB, guid string) error {
	log.Printf("[+] dbase connected to database (preparing insert)")

	res, err := db.Exec("insert into posted (guid, posted) values ($1, $2)",
		guid, true)
	if err != nil {
		log.Printf("[!] dbase couldn't insert into database",
			guid)
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		log.Printf("[!] dbase insert affects 0 rows")
		return fmt.Errorf("insert affects 0 rows")
	}

	return nil
}

func CountStories(db *sql.DB) int64 {
	rows, err := db.Query("select count(*) from posted")
	if err != nil {
		log.Println("[!] dbase select count failed")
		return 0
	}

	var count int64 = 0
	for rows.Next() {
		rows.Scan(&count)
	}

	return count
}
