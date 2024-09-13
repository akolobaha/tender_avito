package domain

import (
	"errors"
	"github.com/google/uuid"
	"tenders/db"
)

type OrganizationResponsible struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id"`
	UserId         string `json:"user_id"`
}

func GetOrganizationResponsible(tenderId uuid.UUID, authorId uuid.UUID) (OrganizationResponsible, error) {
	orgResponsible := OrganizationResponsible{}
	db := db.GetConnection()
	defer db.Close()

	err := db.
		QueryRow(`SELECT id, organization_id, user_id FROM organization_responsible 
							WHERE user_id = $1 and id = 
							   (SELECT organization_responsible_id FROM tender WHERE id = $2)`, authorId, tenderId).
		Scan(&orgResponsible.Id, &orgResponsible.OrganizationId, &orgResponsible.UserId)

	if orgResponsible.Id == "" {
		return OrganizationResponsible{}, errors.New("Указанного ответственного в указанной организации не существует")
	}

	if err != nil {
		return OrganizationResponsible{}, err
	}

	return orgResponsible, nil
}

func IsUserResponsibleToTenderByUsername(username string, tenderId uuid.UUID) error {
	db := db.GetConnection()
	defer db.Close()

	var result bool
	db.QueryRow(`SELECT EXISTS(SELECT *
              FROM organization_responsible
              WHERE user_id = (select id from employee where username = $1)
                AND id = (select organization_responsible_id from tender where id = $2))`, username, tenderId).Scan(&result)

	if !result {
		return errors.New("Не достаточно прав")
	}
	return nil
}

func IsUserResponsibleToTenderByAuthorId(authorId uuid.UUID, tenderId uuid.UUID) error {
	db := db.GetConnection()
	defer db.Close()

	var result bool
	db.QueryRow(`SELECT EXISTS(SELECT *
              FROM organization_responsible
              WHERE user_id = $1
                AND id = (select organization_responsible_id from tender where id = $2))`, authorId, tenderId).Scan(&result)

	if !result {
		return errors.New("Не достаточно прав")
	}
	return nil
}
