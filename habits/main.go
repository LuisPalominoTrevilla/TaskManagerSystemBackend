package main

import (
	"net/http"
	"os"
	"bytes"

	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/db"
	"github.com/joho/godotenv"
	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/routers"
)

func main() {
	godotenv.Load()

	database := db.InitDb()

	r := routers.GetRouter(database)

	var jsonStr = []byte(`{"port":"4001", "service":"habits", "healthCheck":"/check"}`)

	req, err := http.NewRequest("POST", os.Getenv("REGISTRY_HOST")+os.Getenv("REGISTRY_ENDPOINT"), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

	http.ListenAndServe(":4001", r)
}
