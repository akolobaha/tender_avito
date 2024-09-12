package domain

type OrganizationResponsible struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id"`
	UserId         string `json:"user_id"`
}
