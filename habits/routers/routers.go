package routers

import (
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// GetRouter returns the router to be served
func GetRouter(database *mongo.Database) *mux.Router {
	r := mux.NewRouter()
	SetAPIRouter(r, database)
	return r
}
