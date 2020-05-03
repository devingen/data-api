package dto

import "github.com/devingen/api-core/model"

type QueryRequest struct {
	Base        string
	Collection  string
	QueryConfig *model.QueryConfig
}

type QueryResponse struct {
	Results []*model.DataModel `json:"results"`
	Meta    *model.Meta        `json:"meta"`
}
