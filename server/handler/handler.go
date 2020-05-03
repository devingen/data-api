package handler

import (
	"github.com/devingen/data-api/controller"
	"github.com/gorilla/mux"
	"net/http"
)

type ServerHandler struct {
	Controller controller.DataController
	Router     *mux.Router
}

func NewHttpServiceHandler(atamaController controller.DataController) ServerHandler {
	handler := ServerHandler{Controller: atamaController}

	handler.Router = mux.NewRouter()
	handler.Router.HandleFunc("/{base}/{collection}/query", handler.Query).Methods(http.MethodPost)

	return handler
}
