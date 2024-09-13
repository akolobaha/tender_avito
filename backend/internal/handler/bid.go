package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"tenders/db"
	"tenders/internal/domain"
)

func getBidsByTenderIdHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "bids list: %s")
}

func bidCreateHandler(w http.ResponseWriter, r *http.Request) {
	var newBidReq domain.Bid
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newBidReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db := db.GetConnection()
	defer db.Close()

	orgResponsible, err := domain.GetOrganizationResponsible(newBidReq.OrganisationId, newBidReq.CreatorUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bid, err := domain.AddBid(newBidReq, orgResponsible)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, bid)
}

func bidsMyHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	db := db.GetConnection()
	defer db.Close()

	var bids []domain.Bid

	rows, err := db.Query("SELECT id, name, description, status, tender_id\nFROM bid\nWHERE organization_responsible_id IN (SELECT id FROM organization_responsible WHERE user_id IN (SELECT id FROM employee where username = $1));", username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bid domain.Bid

		err = rows.Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bids = append(bids, bid)
	}

	renderJSON(w, bids)

	//fmt.Fprintf(w, "Username: %s", username)
}

func BidsEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["bidId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}

func bidRollback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["bidId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}
