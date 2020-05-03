package handler

import (
	"encoding/json"
	"github.com/devingen/api-core/model"
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/util"
	"github.com/devingen/data-api/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (handler ServerHandler) Query(w http.ResponseWriter, r *http.Request) {

	pathVariables := mux.Vars(r)

	var config model.QueryConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := handler.Controller.Query(&dto.QueryRequest{
		Base:        pathVariables["base"],
		Collection:  pathVariables["collection"],
		QueryConfig: &config,
	})
	response, err := util.BuildResponse(http.StatusOK, result, err)
	server.ReturnResponse(w, response, err)
}
