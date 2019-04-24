package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Controller states all common methods for a controller
type Controller interface {
	Get(http.ResponseWriter, *http.Request)
	initializeController(*mux.Router)
}
