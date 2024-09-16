package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
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
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	bidID := vars["bidId"]
	bidIdUuid, err := uuid.Parse(bidID)

	if err != nil {
		http.Error(w, "Не достаточно прав", http.StatusForbidden)
		return
	}

	err = domain.IsUserResponsibleToBidByUsername(username, bidIdUuid)
	if err != nil {
		http.Error(w, "Не достаточно прав", http.StatusForbidden)
		return
	}

	var newBidPatch domain.Bid
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newBidPatch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := db.GetConnection()
	defer db.Close()

	query := "UPDATE bid SET "
	var args []interface{}
	setClauses := []string{}

	if newBidPatch.Name != "" {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, newBidPatch.Name)
	}

	if newBidPatch.Description != "" {
		setClauses = append(setClauses, "description = $"+strconv.Itoa(len(args)+1))
		args = append(args, newBidPatch.Description)
	}

	if len(setClauses) == 0 {
		http.Error(w, "Нет данных для обновления", http.StatusBadRequest)
		return
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, bidIdUuid)

	_, err = db.Exec(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(`SELECT id, description, status, tender_id, organization_responsible_id, name FROM bid WHERE id = $1`, bidIdUuid).
		Scan(&newBidPatch.Id, &newBidPatch.Description, &newBidPatch.Status, &newBidPatch.TenderId, &newBidPatch.CreatorUsername, &newBidPatch.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, newBidPatch)
}
