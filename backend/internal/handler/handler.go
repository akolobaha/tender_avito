package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/api/ping", ping).Methods("GET")

	// Тендеры
	r.HandleFunc("/api/tenders", tendersList).Methods("GET")
	r.HandleFunc("/api/tenders/new", tenderCreate).Methods("POST")
	r.HandleFunc("/api/tenders/my", tendersByUser).Methods("GET")
	r.HandleFunc("/api/tenders/{tenderId}/edit", tenderUpdate).Methods("PATCH")
	r.HandleFunc("/api/tenders/{tenderId}/rollback/{version}", tenderRollback).Methods("PUT")

	// Предложения
	r.HandleFunc("/api/bids/{tenderId}/list", getBidsByTenderId).Methods("GET")
	r.HandleFunc("/api/bids/new", bidCreate).Methods("POST")
	r.HandleFunc("/api/bids/my", getBidsByUser).Methods("GET")
	r.HandleFunc("/api/bids/{bidId}/edit", bidUpdate).Methods("PATCH")
	r.HandleFunc("api/bids/{bidId}/rollback/{version}", bidRollback).Methods("PUT")

	//slog.Info("Starting server on " + cfg.ServerAddress)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// handler function for the root route
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

// handler function for a custom route
func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "good bye!")

}
