package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"nfip-community-book/data"
	"nfip-community-book/handlers"
)

func main() {
	l := log.New(os.Stdout, "NFIP Community Book: ", log.LstdFlags)
	cb, err := data.GetNFIPCommunityStatusBook(l)

	if err != nil {
		l.Println(err.Error())
		os.Exit(1)
	}

	sh := handlers.NewStatus(l, cb)
	sm := http.NewServeMux()
	sm.Handle("/status", sh)

	s := http.Server{
		Addr:         ":9001",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9001")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	l.Println("Received signal:", sig)

	// wait 30 seconds for existing requests to be completed before shutting down
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
