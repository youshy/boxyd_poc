package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// Choose the folder to serve
	staticDir := "/boxes/"

	// TODO: Figure out how to print the dir
	/*
		router.
			Methods(http.MethodGet).
			Handler(basicAuth(printIndex(staticDir)))
	*/

	// Create the route
	router.
		PathPrefix("/box/").
		Handler(basicAuth(http.StripPrefix("/box/", http.FileServer(http.Dir("."+staticDir)))))

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

const titlepage = `
<html>
{{ range $i := .Body}}
<a href={{$i}}>{{$i}}<br/></a>
{{end}}
</html>
`

type tp struct {
	Title string
	Body  []string
}

func printIndex(p string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var member []string
		path := "." + p

		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if !info.IsDir() {
				temp := strings.Replace(path, "boxes", "box", -1)
				fmt.Printf("replaced from %s to %s\n", path, temp)
				member = append(member, temp)
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		as := tp{Title: "Hello", Body: member}
		t := template.Must(template.New("Index").Parse(titlepage))
		t.Execute(w, as)
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
