package handler

import (
	"encoding/json"
	"github.com/devingen/api-core/server"
	"github.com/devingen/api-core/util"
	"github.com/devingen/data-api/dto"
	"github.com/gorilla/mux"
	"net/http"
)

func (handler ServerHandler) Update(w http.ResponseWriter, r *http.Request) {

	pathVariables := mux.Vars(r)

	var config dto.UpdateConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := handler.Controller.Update(&dto.UpdateRequest{
		Base:                pathVariables["base"],
		Collection:          pathVariables["collection"],
		ID:                  pathVariables["id"],
		UpdateConfig:        &config,
		AuthorizationHeader: r.Header.Get("Authorization"),
	})
	response, err := util.BuildResponse(http.StatusOK, result, err)
	server.ReturnResponse(w, response, err)
}
