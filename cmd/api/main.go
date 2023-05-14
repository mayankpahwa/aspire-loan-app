package main

import (
	"log"

	"github.com/mayankpahwa/aspire-loan-app/app/server"
)

func main() {
	svr, err := server.New()
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
	svr.ListenAndServe()
}
