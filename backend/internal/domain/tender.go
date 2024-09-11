package domain

import (
	_ "github.com/lib/pq"
)

type TenderResp struct {
	ID             int    `json:"id"`
	Description    string `json:"description"`
	OrganizationId string `json:"organization_id"`
	EmployeeId     string `json:"employee_id"`
	Version        string `json:"version_uid"`
	Status         string `json:"status"`
}
