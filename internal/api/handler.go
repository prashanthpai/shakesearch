package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Searcher interface {
	Search(query string, filter string) []string
	Filters() []string
}

type Handler struct {
	searcher Searcher
}

func NewHandler(s Searcher) http.Handler {
	h := &Handler{
		searcher: s,
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/search", h.search)
	mux.HandleFunc("/filters", h.listFilters)

	return mux
}

func (h *Handler) listFilters(w http.ResponseWriter, r *http.Request) {
	filters := h.searcher.Filters()

	b, err := json.Marshal(filters)
	if err != nil {
		log.Printf("json.Marshal() failed: %s", err.Error())
		http.Error(w, "json encoding failure", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Printf("w.Write() failed: %s", err.Error())
	}
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	query, ok := r.URL.Query()["q"]
	if !ok || len(query[0]) < 1 {
		http.Error(w, "missing search query in URL params", http.StatusBadRequest)
		return
	}

	var filter string
	if q, ok := r.URL.Query()["filter"]; ok {
		if len(q[0]) != 0 {
			filter = q[0]
		}
	}

	results := h.searcher.Search(query[0], filter)

	b, err := json.Marshal(results)
	if err != nil {
		log.Printf("json.Marshal() failed: %s", err.Error())
		http.Error(w, "json encoding failure", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Printf("w.Write() failed: %s", err.Error())
	}
}
