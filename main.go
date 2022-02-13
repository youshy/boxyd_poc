package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// Choose the folder to serve
	staticDir := "/boxes/"

	// Create the route
	router.
		PathPrefix(staticDir).
		Handler(basicAuth(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir)))))

	log.Println("Available routes:")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		fmt.Printf("%s\t%s\n", m, t)
		return nil
	})

	log.Println("Starting server...")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		usernameCheck := os.Getenv("BOXYD_USERNAME")
		passwordCheck := os.Getenv("BOXYD_PASSWORD")
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(usernameCheck))
			expectedPasswordHash := sha256.Sum256([]byte(passwordCheck))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
