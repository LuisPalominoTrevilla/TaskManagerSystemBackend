package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HabitsController holds important data for the habits controller
type HabitsController struct {
	data string
}

// Get serves as a GET request
func (controller *HabitsController) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Habbits controller")
}

func (controller *HabitsController) initializeController(r *mux.Router) {
	r.HandleFunc("/", controller.Get).Methods(http.MethodGet)
}

// SetController sets the controller for the sets up the habits controllet
func SetController(r *mux.Router) {
	habitsController := HabitsController{data: "nombre"}
	habitsController.initializeController(r)
}
