package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sheenobu/go-examples/rest/users"
)

func main() {
	r := httprouter.New()

	// build the user system
	userStorage := users.NewStorage()
	users.RegisterHTTP(r, userStorage)

	//TODO: build the post system

	log.Printf("Listening on http://127.0.0.1:8080/")
	log.Fatal(http.ListenAndServe(":8080", r))
}
