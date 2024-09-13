package domain

import (
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const TenderStatusCreated = "CREATED"
const TenderStatusPublished = "PUBLISHED"
const TenderStatusClosed = "CLOSED"
const TenderStatusOpen = "OPEN"

const ServiceTypeConstruction = "Construction"
const ServiceTypeDelivery = "Delivery"
const ServiceTypeManufacture = "Manufacture"

type TenderResp struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ServiceType    string `json:"serviceType"`
	Status         string `json:"status"`
	OrganizationId string `json:"organizationId"`
	EmployeeId     string `json:"employeeId"`
	Version        string `json:"version"`
}

type Tender struct {
	Id              uuid.UUID `json:"id,omitempty"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"serviceType"`
	Status          string    `json:"status"`
	OrganizationId  uuid.UUID `json:"organizationId"`
	CreatorUserName string    `json:"creatorUserName"`
}

type TenderPatch struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	ServiceType string    `json:"serviceType,omitempty"`
}
