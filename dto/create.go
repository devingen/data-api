package dto

type CreateRequest struct {
	Base                string
	Collection          string
	CreateConfig        *CreateConfig
	AuthorizationHeader string
}

type CreateConfig struct {
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}

type CreateResponse struct {
	ID       string `json:"_id"`
	Revision int    `json:"_revision"`
}
