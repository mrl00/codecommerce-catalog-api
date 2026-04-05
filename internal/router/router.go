package router

import (
	"goapi/internal/handler"

	"github.com/gorilla/mux"
)

func New() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/health", handler.Health).Methods("GET")

	return r
}
