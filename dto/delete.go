package dto

type DeleteRequest struct {
	Base                string
	Collection          string
	ID                  string
	AuthorizationHeader string
}

type DeleteResponse struct{}
