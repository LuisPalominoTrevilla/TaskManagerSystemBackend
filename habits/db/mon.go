package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// Model is just a wrapper for a model
type Model interface{}

// Database represents an interface to the basic oprations of a model
type Database interface {
	Get(bson.D, *Model)
	Insert(Model)
}

type key string

const (
	hostKey     = key("hostKey")
	usernameKey = key("usernameKey")
	passwordKey = key("passwordKey")
	databaseKey = key("databaseKey")
)

// InitDb initializes the database
func InitDb() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, hostKey, os.Getenv("MONGO_HOST")+os.Getenv("MONGO_PORT"))
	ctx = context.WithValue(ctx, usernameKey, os.Getenv("MONGO_USER"))
	ctx = context.WithValue(ctx, passwordKey, os.Getenv("MONGO_PASS"))
	ctx = context.WithValue(ctx, databaseKey, os.Getenv("MONGO_DB"))

	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		ctx.Value(usernameKey).(string),
		ctx.Value(passwordKey).(string),
		ctx.Value(hostKey).(string),
		ctx.Value(databaseKey).(string),
	)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.NewClient(clientOptions)
	
	if err != nil {
		println(err.Error())
		log.Fatal("There was an error creating a new client")
	}

	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		println(err.Error())
		log.Fatal("There was an error connecting to the database")
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	primary := readpref.Primary()
	err = client.Ping(ctx, primary)

	if err != nil {
		log.Fatal("Couldn't find a server ", err)
	}
	fmt.Println("Connected to MongoDB!")
	return client.Database(os.Getenv("MONGO_DB"))
}