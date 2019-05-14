package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Habit is used to create a habit
type Habit struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title      string `json:"title"`
	Type	   int 	  `json:"type"`
	Difficulty int 	  `json:"difficulty"`
	UserId     string `json:"userId"`
	Image	   string `json:"image"`
	Score	   float32	  `json:"score"`
}
