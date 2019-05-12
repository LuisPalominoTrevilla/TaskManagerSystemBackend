package main

import (
	"net/http"
	"fmt"

	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/db"
	"github.com/joho/godotenv"
	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/routers"
)

func main() {

	godotenv.Load()

	fmt.Println("Loaded .env successfully")

	database := db.InitDb()

	r := routers.GetRouter(database)
	http.ListenAndServe(":4001", r)
}
