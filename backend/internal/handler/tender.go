package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func tendersList(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Tenders list: %s")
}

func tenderCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tenders create: %s")
}

func tendersByUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Username: %s", username)
}

func tenderUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["tenderId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}

func tenderRollback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["tenderId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}
