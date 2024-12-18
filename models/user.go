package models

import "time"

type User struct {
	Id        string    `json:"user_id" bson:"user_id"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"password" bson:"password"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	RemovedAt time.Time `json:"removed_at,omitempty" bson:"removed_at,omitempty"`
}

type UserDelete struct {
	UserId string `json:"user_id" bson:"user_id"`
}

type UserLogin struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
