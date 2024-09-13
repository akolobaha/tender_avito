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
	r.HandleFunc("/api/tenders", tendersListHandler).Methods("GET")
	r.HandleFunc("/api/tenders/new", tenderCreateHandler).Methods("POST")
	r.HandleFunc("/api/tenders/my", tendersMyHandler).Methods("GET")
	r.HandleFunc("/api/tenders/{tenderId}/edit", tenderUpdateHandler).Methods("PATCH")
	r.HandleFunc("/api/tenders/{tenderId}/rollback/{version}", tenderRollback).Methods("PUT")

	// Предложения
	r.HandleFunc("/api/bids/{tenderId}/list", getBidsByTenderIdHandler).Methods("GET")
	r.HandleFunc("/api/bids/new", bidCreateHandler).Methods("POST")
	r.HandleFunc("/api/bids/my", bidsMyHandler).Methods("GET")
	r.HandleFunc("/api/bids/{bidId}/edit", BidsEditHandler).Methods("PATCH")
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

func parseOffsetParams(r *http.Request) (string, string) {
	// Извлечение параметров limit, offset и serviceType из URL
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	// Установка значений по умолчанию, если не указаны
	if limit == "" {
		limit = "5" // по умолчанию 5
	}
	if offset == "" {
		offset = "0" // По умолчанию 0
	}

	return limit, offset
}
