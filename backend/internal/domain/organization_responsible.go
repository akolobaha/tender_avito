package domain

import (
	"errors"
	"tenders/db"
)

type OrganizationResponsible struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id"`
	UserId         string `json:"user_id"`
}

func GetOrganizationResponsible(organizationId string, username string) (OrganizationResponsible, error) {
	orgResponsible := OrganizationResponsible{}
	db := db.GetConnection()
	defer db.Close()

	err := db.
		QueryRow("SELECT id FROM organization_responsible WHERE user_id = (select id from employee where username = $1) and organization_id = $2", username, organizationId).
		Scan(&orgResponsible.Id)

	if orgResponsible.Id == "" {
		return OrganizationResponsible{}, errors.New("Указанного ответственного в указанной организации не существует")
	}

	if err != nil {
		return OrganizationResponsible{}, err
	}

	return orgResponsible, nil
}
