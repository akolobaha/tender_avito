package domain

import (
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TenderResp struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
	EmployeeId     string `json:"employeeId"`
	Version        string `json:"version"`
	Status         string `json:"status"`
}

type TenderReq struct {
	Id              uuid.UUID `json:"id,omitempty"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"serviceType"`
	Status          string    `json:"status"`
	OrganizationId  uuid.UUID `json:"organizationId"`
	CreatorUserName string    `json:"creatorUserName"`
}
