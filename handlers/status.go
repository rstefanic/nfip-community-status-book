package handlers

import (
	"log"
	"net/http"

	"nfip-community-book/data"
)

type Status struct {
	l  *log.Logger
	cb data.NFIPCommunityStatuses
}

func NewStatus(l *log.Logger, cb data.NFIPCommunityStatuses) Status {
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
	search := queries.Get("search")

	if len(search) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	s.l.Printf("[STATUS] Requested search for term \"%s\"\n", search)
	communities := s.cb.Search(search)
	err := communities.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
