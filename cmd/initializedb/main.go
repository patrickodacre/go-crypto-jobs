package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func main() {

	var dbPath = flag.String("dbpath", "my.db", "Path to Database")

	db, err := bolt.Open(*dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Jobs"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

}
