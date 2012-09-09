// bot contains the portion of the code responsible for actually posting stories
// and updating the database.
package bot

import (
	"database/sql"
	"fmt"
	"github.com/gokyle/golobsters/dbase"
	"github.com/gokyle/twitter"
	rss "github.com/jteeuwen/go-pkg-rss"
	"log"
	"os"
	"strings"
	"time"
)

// urls
var feedUri = os.Getenv("RSS_FEED")

// 140 characters - length of a t.co link
const maxTwitterStatus = 115
const maxADNStatus = 256

// update variables
var lastUpdate time.Time

var numWorkers = 3
var newStories = make(chan story, 5)
var twitterApi twitter.Twitter

type story struct {
	title string
	guid  string
	link  string
}

func LastUpdate() string {
	noTime := new(time.Time)
	if lastUpdate == *noTime {
		return ""
	}
	return lastUpdate.String()
}

func Story(item *rss.Item) story {
	s := story{item.Title, item.Guid, item.Links[0].Href}
	return s
}

// Status returns a status message truncated to the requested length
func Status(message string, length int) string {
	if len(message) < length {
		return message
	}

	words := strings.Fields(message)
	status := ""
	for _, word := range words {
		if len(status)+len(word)+1 < (len(status) - 3) {
			status = status + " " + word
		} else {
			break
		}
	}
	status += "..."
	return status
}

// TwitterStatus returns an appropriate status for a Twitter status update 
// given a title and link.
func TwitterStatus(title string, link string) string {
	return fmt.Sprintf("%s %s", Status(title, maxTwitterStatus), link)
}

// ADNStatus returns an appropriate status for an App Dot Net status update 
// given a title and link.
func ADNStatus(title string, link string) string {
	// ADN doesn't use automatic URL-shortening like twitter
	length := maxADNStatus - len(link) - 1
	return fmt.Sprintf("%s %s", Status(title, length), link)
}

// Given an RSS feed item, determine whether it exists in the database and
// if not, post it. This is designed such that it can be run from a
// goroutine.
func (s story) process(db *sql.DB) error {
	if posted, err := dbase.StoryPosted(db, s.guid); err != nil {
		log.Printf("[!] bot StoryHandler failure: %s\n", err)
		return err
	} else if posted {
		log.Printf("[+] bot skipping %s, already posted\n", s.guid)
		return nil
	}

	// story hasn't been posted
	log.Printf("[+] bot worker posting story\n")
	if err := s.post(); err != nil {
		log.Printf("[!] error posting status: %s\n", err)
		return err
	} else if err = dbase.PostStory(db, s.guid); err != nil {
		// once we've posted to twitter, we need to make sure
		// the database is updated!
		var errors int64 = 1
		for {
			log.Printf("[!] %d errors posting to database", errors)
			if err = dbase.PostStory(db, s.guid); err != nil {
				break
			}
			errors++
			time.Sleep(1)
		}
	}

	log.Println("[+] bot successful update")
	return nil
}

// PostStatus is responsible for actually posting the story. It assumes the
// story has not already been posted (otherwise an error will be returned).
// A nil return means the appropriate action for the story has been taken,
// whether skipping over it or updating the database.
func (s story) post() (err error) {
	status := TwitterStatus(s.title, s.guid)
	_, err = twitterApi.Tweet(status)
	log.Println("[-] err: ", err)
	return err
}

func getStories() error {
	timeout := 5          // 5 seconds
	feedTarget := feedUri // rss feed to follow
	feed := rss.New(timeout, true, nil, txNewItems)
	for {
		if err := feed.Fetch(feedTarget, nil); err != nil {
			log.Printf("bot feed failure %s: %s", feedTarget, err)
			return err
		}

		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}

	return nil
}

// Kick off the bot with Run(). Its signature matches the one required by
// gomon. When Run() is called, the environment should already be set up.
func Run() error {
	log.Println("[+] bot starts")

	log.Println("[+] bot initialising twitter API connection")
	twitterApi = twitter.Twitter{
		ConsumerKey:      os.Getenv("TW_CKEY"),
		ConsumerSecret:   os.Getenv("TW_CSEC"),
		OAuthToken:       os.Getenv("TW_ATOK"),
		OAuthTokenSecret: os.Getenv("TW_ASEC"),
	}

	log.Println("[+] bot starting worker pool")
	for i := 0; i < numWorkers; i++ {
		go worker(int8(i))
	}

	log.Println("[+] bot starting feed monitor")
	err := getStories()

	return err
}

func worker(id int8) {
	db, err := dbase.ConnectFromEnv()
	if err != nil {
		log.Println("[+] bot connected to database (preparing select)")
	}
	defer db.Close()
	for {
		s := <-newStories
		err := s.process(db)
		if err != nil {
			log.Printf("[!] worker %d error processing story: %s",
				id, err)
		}
	}
        log.Printf("[!] worker %d dies!\n", id)
}

func txNewItems(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	log.Printf("[+] bot %d new stories on %s\n", len(newitems), feed.Url)
	lastUpdate = time.Now()
	for _, item := range newitems {
		newStories <- Story(item)
	}
}
