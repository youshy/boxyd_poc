package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/boxyd/qrcoder"
	"github.com/gorilla/mux"
)

func handleSingleItem(dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boxID, err := strconv.Atoi(vars["box_id"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		p := fmt.Sprintf(".%s/%04d.html", dir, boxID)
		log.Printf("path: %v\n", p)

		http.ServeFile(w, r, p)
	})
}

func generateQR() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		boxID, err := strconv.Atoi(vars["box_id"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		path := fmt.Sprintf("%s%s", r.Host, r.URL.Path)
		log.Printf("path: %s\tboxID: %v\n", path, boxID)

		label, err := qrcoder.Generate(path, "BOXYD", boxID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(label)
	})
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
