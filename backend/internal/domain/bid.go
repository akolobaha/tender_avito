package domain

import (
	"github.com/google/uuid"
	"tenders/db"
)

const StatusCreated = "CREATED"

type Bid struct {
	Id              uuid.UUID `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	TenderId        uuid.UUID `json:"tenderId"`
	OrganisationId  string    `json:"organizationId,omitempty"`
	CreatorUsername string    `json:"creatorUsername,omitempty"`
}

type BidReq struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TenderId    uuid.UUID `json:"tenderId"`
	AuthorType  string    `json:"authorType"`
	AuthorId    uuid.UUID `json:"authorId"`
}

func AddBid(bid BidReq, orgResponsible OrganizationResponsible) (BidReq, error) {
	db := db.GetConnection()
	defer db.Close()

	err := db.QueryRow(`INSERT INTO bid 
    	(name, description, status, tender_id, organization_responsible_id, author_type) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		bid.Name, bid.Description, StatusCreated, bid.TenderId, orgResponsible.Id, bid.AuthorType).Scan(&bid.Id)

	return bid, err
}
