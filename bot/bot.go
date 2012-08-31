// bot contains the portion of the code responsible for actually posting stories
// and updating the database.
package bot

import (
        "fmt"
        "golobsters/lobsterdb"
        "time"
        "log"
)

// 140 characters - length of a t.co link
const maxTwitterStatus = 115
const maxADNStatus = 256

type story struct {
        title string
        guid  string
}

// Status returns a status message truncated to the requested length
func Status(message string, length int) string {
        if len(message) < length {
                return message
        }

        words := message.Fields()
        status := ""
        for _, word := range words {
                if len(status) + len(word) + 1 < (message - 3) {
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
func TwitterStatus(title string, link string) {
        return fmt.Sprintf("%s %s", Status(title, maxTwitterStatus), link)
}

// ADNStatus returns an appropriate status for an App Dot Net status update 
// given a title and link.
func ADNStatus(title string, link string) {
        // ADN doesn't use automatic URL-shortening like twitter
        length := maxADNStatus - len(link) - 1
        return fmt.Sprintf("%s %s", Status(title, length), link)
}

// Given an RSS feed item, determine whether it exists in the database and
// if not, post it. This is designed such that it can be run from a
// goroutine.
func (s story) process() error {
        if posted, err := lobsterdb.StoryPosted(s.guid); err != nil {
                log.Printf("[!] bot StoryHandler failure: %s\n", err)
                return err
        } else if posted {
                log.Printf("[+] bot skipping %s, already posted\n")
                return nil
        }

        // story hasn't been posted
        if err := s.post(); err != nil {
                log.Printf("[!] error posting status: %s\n", err)
                return err
        } else if err = lobsterdb.PostStory(s.guid); err != nil {
                // once we've posted to twitter, we need to make sure
                // the database is updated!
                errors := 1.(int64)
                for {
                        log.Printf("[!] %d errors posting to database", errors)
                        if err = lobsterdb.PostStory(s.guid); err != nil {
                                break
                        }
                        errors++
                        time.Sleep(1)
                }
        }

        return nil
}

// PostStatus is responsible for actually posting the story. It assumes the
// story has not already been posted (otherwise an error will be returned).
// A nil return means the appropriate action for the story has been taken,
// whether skipping over it or updating the database.
func (s story) post (err error) {
        //status := TwitterStatus(s.title, s.guid)

        return fmt.Errorf("status updates aren't implemented")
}


// Update retrieves the top five stories from the lobste.rs homepage, and
// passes the last five to the StoryHandler.
func Update() error {
        
}
