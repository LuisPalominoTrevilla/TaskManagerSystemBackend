package routers

import (
	"fmt"
	"net/http"

	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/controllers"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/gorilla/mux"
)

func handleAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API is alive")
}

// SetAPIRouter sets the API Router
func SetAPIRouter(r *mux.Router, db *mongo.Database) {
	habitRouter := r.PathPrefix("/habits").Subrouter()
	habitRouter.StrictSlash(false)
	controllers.SetHabitsController(habitRouter, db)
}