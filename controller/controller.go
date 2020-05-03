package controller

import "github.com/devingen/data-api/dto"

type DataController interface {
	Query(request *dto.QueryRequest) (dto.QueryResponse, error)
}
