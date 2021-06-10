package models

type User struct {
	ID         uint   `json:"id"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Password   string `json:"password"`
	Email      string `json:"email"`
}

// Users struct
type Users struct {
	Users []User `json:"Users"`
}

type LoginUser struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
