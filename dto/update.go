package dto

type UpdateRequest struct {
	Base                string
	Collection          string
	ID                  string
	UpdateConfig        *UpdateConfig
	AuthorizationHeader string
}

type UpdateConfig struct {
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}

type UpdateResponse struct {
	Updated  string `json:"_updated"`
	Revision int    `json:"_revision"`
}
