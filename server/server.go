package server

import (
	"log"
	"net/http"
	"shivamaravanthe/HosangadiReports/api"
	"shivamaravanthe/HosangadiReports/constants"

	"github.com/gorilla/mux"
)

func Server() {
	router := mux.NewRouter()

	router.HandleFunc("/gst", api.Gst).Methods("GET")

	server := &http.Server{
		Addr:    constants.PORT,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to connect to database %v", err)
		return
	}
}
