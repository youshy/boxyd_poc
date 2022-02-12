package main

import (
	"net/http"
)

func main() {
	http.Handle("/boxes/", http.StripPrefix("/boxes/", http.FileServer(http.Dir("./pages"))))
	http.ListenAndServe(":3000", nil)
}
