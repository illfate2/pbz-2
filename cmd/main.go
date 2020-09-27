package main

import (
	"log"
	"net/http"

	"pbz2/pkg/api"
)

func main() {
	s := api.NewServer()
	log.Fatal(http.ListenAndServe(":8080", s))
}
