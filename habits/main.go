package main

import (
	"net/http"

	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/db"
	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/routers"
)

func main() {

	database := db.InitDb()

	r := routers.GetRouter(database)
	http.ListenAndServe(":4001", r)
}
