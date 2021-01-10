package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prashanthpai/shakesearch/pkg/shake"
)

const (
	inputFile = "completeworks.txt"
)

func main() {
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("os.Open(%s) failed: %s", inputFile, err.Error())
	}
	defer f.Close()

	shakeSearcher := shake.NewSearcher()

	if err := shakeSearcher.Load(f); err != nil {
		log.Fatalf("shakeSearcher.Load() failed: %s", err.Error())
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(shakeSearcher))
	http.HandleFunc("/filters", handleListFilters(shakeSearcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher interface {
	Search(query string, filter string) []string
	Filters() []string
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("missing search query in URL params"))
			return
		}

		var filter string
		if q, ok := r.URL.Query()["filter"]; ok {
			if len(q[0]) != 0 {
				filter = q[0]
			}
		}

		results := searcher.Search(query[0], filter)

		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(results); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("encoding failure"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(buf.Bytes())
	}
}

func handleListFilters(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		results := searcher.Filters()

		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(results); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("encoding failure"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(buf.Bytes())
	}
}
