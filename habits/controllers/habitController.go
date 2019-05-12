package controllers

import (
	"fmt"
	"net/http"

	database "github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/db"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/gorilla/mux"
)

// HabitsController holds important data for the habits controller
type HabitsController struct {
	habitsDB *database.HabitsDB
}

// Get serves as a GET request
func (controller *HabitsController) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Habits controller")
}

func (controller *HabitsController) initializeController(r *mux.Router) {
	r.HandleFunc("/", controller.Get).Methods(http.MethodGet)
}

// SetHabitsController sets the controller for the sets up the habits controllet
func SetHabitsController(r *mux.Router, db *mongo.Database) {
	habit := database.HabitsDB{Habits: db.Collection("habits")}
	
	habitsController := HabitsController{habitsDB: &habit}
	habitsController.initializeController(r)
}
