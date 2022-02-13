package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	box := "/box"

	router.Handle(box+"/{box_id}", basicAuth(getInfo(handleSingleItem(box)))).Methods(http.MethodGet)

	router.Handle(box+"/{box_id}/qr", getInfo(generateQR(false))).Methods(http.MethodGet)
	router.Handle(box+"/{box_id}/qrsmall", getInfo(generateQR(true))).Methods(http.MethodGet)

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
