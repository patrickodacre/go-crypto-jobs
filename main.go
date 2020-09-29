package main

import (
	// "encoding/json"
	"flag"
	// "fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	// "github.com/go-kit/kit/endpoint"
	// httptransport "github.com/go-kit/kit/transport/http"
)

func main() {

	var (
		domain = flag.String("domain", ":8080", "Top Level Domain")
	)

	flag.Parse()

	r := chi.NewRouter()

	scrape()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
		Addr:         *domain,
	}

	log.Fatal(srv.ListenAndServe())
}
