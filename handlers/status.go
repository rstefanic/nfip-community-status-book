package handlers

import (
	"log"
	"net/http"

	"nfip-community-book/data"
)

type Status struct {
	l  *log.Logger
	cb data.NFIPCommunities
}

func NewStatus(l *log.Logger, cb data.NFIPCommunities) Status {
	return Status{l, cb}
}

func (s Status) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.getStatus(rw, r)
		return
	}

	rw.WriteHeader(http.StatusBadRequest)
}

func (s Status) getStatus(rw http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	term := queries.Get("search")

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
