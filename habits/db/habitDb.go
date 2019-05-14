package db
import (
	"context"
	"fmt"
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
func (db *HabitsDB) Get(filter bson.D) ([]*models.Habit, error) {
	cur, err := db.Habits.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result []*models.Habit

	for cur.Next(context.TODO()) {
		var habit models.Habit
		err := cur.Decode(&habit)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		result = append(result, &habit)
	}
	return result, nil
}
// UpdateOne updates a document in the database
func (db *HabitsDB) UpdateOne(filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	return db.Habits.UpdateOne(context.TODO(), filter, update)
}

// DeleteOne updates a document in the database
func (db *HabitsDB) DeleteOne(filter bson.D) (*mongo.DeleteResult, error) {
	return db.Habits.DeleteOne(context.TODO(), filter)
}