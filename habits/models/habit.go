package models

import "github.com/mongodb/mongo-go-driver/bson/primitive"

//Habit is used to create a habit
type Habit struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title      string `json:"title"`
	Type	   int 	  `json:"type"`
	Difficulty int 	  `json:"difficulty"`
	UserEmail  string `json:"userEmail"`
	Image	   string `json:"image"`
	Score	   int	  `json:"score"`
}
