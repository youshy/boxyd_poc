package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

	log.Println("Starting server...")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
