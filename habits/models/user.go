package models

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}
