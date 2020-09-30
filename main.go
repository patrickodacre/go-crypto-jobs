package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi"
	// "github.com/go-kit/kit/endpoint"
	// httptransport "github.com/go-kit/kit/transport/http"
	"github.com/boltdb/bolt"
	"github.com/patrickodacre/go-crypto-jobs/cryptojobslist"
)

func main() {

	var (
		domain = flag.String("domain", ":8080", "Top Level Domain")
	)

	flag.Parse()

	// connect to our db
	var db *bolt.DB
	{
		database, err := bolt.Open("my.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		db = database
		defer db.Close()
	}

	r := chi.NewRouter()

	// routes:
	{
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World"))
		})

		r.Get("/records/jobs", func(w http.ResponseWriter, r *http.Request) {

			numOfRecords := 0
			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Jobs"))

				numOfRecords = b.Stats().KeyN
				return nil
			})

			w.Write([]byte(fmt.Sprintf("You have %d jobs stored.", numOfRecords)))
		})

		r.Get("/crypto-jobs-list", func(w http.ResponseWriter, r *http.Request) {

			jobs := []cryptojobslist.Job{}

			filters, hasFilters := r.URL.Query()["filter"]

			// some terms like "go" will return a lot of false positives
			filterTerm := filters[0]
			{
				switch {
				case filterTerm == "go":
					filterTerm = "golang"
				}
			}

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Jobs"))

				b.ForEach(func(k, v []byte) error {
					var j cryptojobslist.Job

					json.Unmarshal(v, &j)

					if !hasFilters {

						jobs = append(jobs, j)
						return nil
					}

					re := regexp.MustCompile(`(?i)` + filterTerm)

					if re.MatchString(j.JobDescription) || re.MatchString(j.Skills) {
						jobs = append(jobs, j)
					}

					return nil
				})

				return nil
			})

			log.Printf("Returing %d jobs", len(jobs))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(jobs)
		})
	}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
		Addr:         *domain,
	}

	log.Fatal(srv.ListenAndServe())
}
