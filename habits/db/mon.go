package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

// Model is just a wrapper for a model
type Model interface{}

// Database represents an interface to the basic oprations of a model
type Database interface {
	Get(bson.D, *Model)
	Insert(Model)
}

// InitDb initializes the database
func InitDb() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, os.Getenv("MONGO_HOST")+":"+os.Getenv("MONGO_PORT"))
	if err != nil {
		log.Fatal("There was an error connecting to the database")
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Couldnt find a server ", err)
	}
	fmt.Println("Connected to MongoDB!")
	return client.Database(os.Getenv("MONGO_DB"))
}
