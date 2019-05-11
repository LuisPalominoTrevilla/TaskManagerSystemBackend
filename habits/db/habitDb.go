package db

import (
	"context"

	"github.com/LuisPalominoTrevilla/TaskManagerSystemBackend/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// UserDB represents the habits collection
type HabitsDB struct {
	Habits *mongo.Collection
}

func (db *HabitsDB) GetByID(filter bson.D, result *models.Habit) error {
	return db.Habits.FindOne(context.TODO(), filter).Decode(&result)
}

func (db *HabitsDB) Insert(habit models.Habit) (*mongo.InsertOneResult, error) {
	return db.Habits.InsertOne(context.TODO(), habit)
}
