package handlers

import (
	"log"
	"net/http"

	"nfip-community-book/data"
)

type Rating struct {
	l   *log.Logger
	crs data.NFIPCommunityRatings
}

func NewRating(l *log.Logger, crs data.NFIPCommunityRatings) Rating {
	return Rating{l, crs}
}

func (rh Rating) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rh.getStatus(rw, r)
		return
	}

	rw.WriteHeader(http.StatusBadRequest)
}

func (rh Rating) getStatus(rw http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	search := queries.Get("search")

	if len(search) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rh.l.Printf("[RATING] Request search for term \"%s\"\n", search)
	communityRatings := rh.crs.Search(search)
	err := communityRatings.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
