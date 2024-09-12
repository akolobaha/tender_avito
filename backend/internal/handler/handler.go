package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"tenders/internal/config"
)

func NewRouter(cfg *config.Config) {
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

	slog.Info("Starting server on " + cfg.ServerAddress)
	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
