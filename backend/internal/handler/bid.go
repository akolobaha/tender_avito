package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func getBidsByTenderId(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "bids list: %s")
}

func bidCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bids create: %s")
}

func getBidsByUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Username: %s", username)
}

func bidUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["bidId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}

func bidRollback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["bidId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}
