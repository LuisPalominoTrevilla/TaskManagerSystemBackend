package main

import (
	"net/http"

	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/db"
	"github.com/joho/godotenv"
	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/routers"
)

func main() {
	godotenv.Load()

	database := db.InitDb()

	r := routers.GetRouter(database)
	http.ListenAndServe(":4001", r)
}
