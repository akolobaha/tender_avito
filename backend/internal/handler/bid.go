package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"tenders/db"
	"tenders/internal/domain"
)

func bidCreateHandler(w http.ResponseWriter, r *http.Request) {
	var newBidReq domain.BidReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newBidReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db := db.GetConnection()
	defer db.Close()

	var orgResponsible domain.OrganizationResponsible
	orgResponsible, err = domain.GetOrganizationResponsible(newBidReq.TenderId, newBidReq.AuthorId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	limit, offset := parseOffsetParams(w, r)

	db := db.GetConnection()
	defer db.Close()

	var bids []domain.Bid

	// Измененный SQL-запрос с LIMIT и OFFSET
	query := `
        SELECT id, name, description, status, tender_id
        FROM bid
        WHERE organization_responsible_id IN (
            SELECT id 
            FROM organization_responsible 
            WHERE user_id IN (
                SELECT id 
                FROM employee 
                WHERE username = $1
            )
        )
        LIMIT $2 OFFSET $3;`

	rows, err := db.Query(query, username, limit, offset)
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
}

func BidsEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["bidId"]
	fmt.Fprintf(w, "User ID: %s", userID)
}
