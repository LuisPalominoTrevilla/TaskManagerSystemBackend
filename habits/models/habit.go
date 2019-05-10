package models

//Habit is used to create a habit
type Habit struct {
	Name       string `json:"name"`
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
	UserEmail  string `json:"userEmail"`
}
