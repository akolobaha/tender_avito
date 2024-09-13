package domain

import "tenders/db"

type Bid struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	TenderId        string `json:"tenderId"`
	OrganisationId  string `json:"organizationId,omitempty"`
	CreatorUsername string `json:"creatorUsername,omitempty"`
}

func AddBid(bid Bid, orgResp OrganizationResponsible) (Bid, error) {
	db := db.GetConnection()
	defer db.Close()

	err := db.QueryRow("INSERT INTO bid (name, description, status, tender_id, organization_responsible_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", bid.Name, bid.Description, "CREATED", bid.TenderId, orgResp.Id).Scan(&bid.Id)

	return bid, err
}
