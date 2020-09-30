package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi"
	// "github.com/go-kit/kit/endpoint"
	// httptransport "github.com/go-kit/kit/transport/http"
	"github.com/patrickodacre/go-crypto-jobs/cryptojobslist"
)

func main() {

	var (
		domain = flag.String("domain", ":8080", "Top Level Domain")
	)

	flag.Parse()

	r := chi.NewRouter()

	jobs := cryptojobslist.Scrape()

	log.Println("scraping done")

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})

	r.Get("/crypto-jobs-list", func(w http.ResponseWriter, r *http.Request) {

		filters, hasFilters := r.URL.Query()["filter"]

		if !hasFilters {

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(jobs)
			return
		}

		// filter jobs
		filteredJobs := []cryptojobslist.Job{}

		for _, j := range jobs {
			re := regexp.MustCompile(`(?i)` + filters[0])

			if re.MatchString(j.JobDescription) || re.MatchString(j.Skills) {
				filteredJobs = append(filteredJobs, j)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filteredJobs)
	})

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
		Addr:         *domain,
	}

	log.Fatal(srv.ListenAndServe())
}
