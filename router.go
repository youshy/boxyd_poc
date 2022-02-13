package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	box := "/box/"
	gear := "/gear/"

	router.
		PathPrefix(gear).
		Handler(basicAuth(handleSingleItem()))

	router.Handle(box, basicAuth(handleSingleItem())).Methods(http.MethodGet)

	router.Handle(box+"{box_id}/qr", basicAuth(generateQR())).Methods(http.MethodGet)

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

	return router
}
