package handlers

import (
	"log"
	"net/http"

	"nfip-community-book/data"
)

type Search struct {
	l  *log.Logger
	cb data.NFIPCommunities
}

func NewSearch(l *log.Logger, cb data.NFIPCommunities) Search {
	return Search{l, cb}
}

func (s Search) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.getSearch(rw, r)
		return
	}

	rw.WriteHeader(http.StatusBadRequest)
}

func (s Search) getSearch(rw http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	term := queries.Get("term")

	if len(term) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.l.Printf("Requested search for term \"%s\"\n", term)
	communities := s.cb.Search(term)
	err := communities.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
