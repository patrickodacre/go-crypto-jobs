package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"log"

	"github.com/boltdb/bolt"
	cryptojobslist "github.com/patrickodacre/go-crypto-jobs/cryptojobslist"
)

func main() {

	var (
		site   = flag.String("site", "crypto-jobs-list", "The site url you want to scrape")
		dbPath = flag.String("dbpath", "my.db", "Path to Database")
	)

	db, err := bolt.Open(*dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	switch {
	case *site == "crypto-jobs-list":
		scrapeCryptoJobsList(db)
	default:
		panic("No site")
	}
}

func scrapeCryptoJobsList(db *bolt.DB) {
	jobs := cryptojobslist.Scrape()

	for _, j := range jobs {

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Jobs"))

			// make a unique key with job id and title
			// to avoid duplicating db records
			key := sha256.Sum256([]byte(j.ID + " " + j.JobTitle))

			jobData, err := json.Marshal(j)

			if err != nil {
				log.Fatalf("Error marshaling scraped jobs", err)
				return err
			}

			err = b.Put(key[:], jobData)

			if err != nil {
				log.Fatalf("Error putting scraped jobs", err)
				return err
			}

			return err
		})
	}
}
