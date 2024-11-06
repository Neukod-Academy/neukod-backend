package models

import "time"

type Student struct {
	StudentId      string          `json:"id" bson:"id"`
	Name           string          `json:"name" bson:"name"`
	Age            uint8           `json:"age" bson:"age"`
	Grade          uint8           `json:"grade" bson:"grade"`
	BornDay        time.Time       `json:"born_day" bson:"born_day"`
	RegisteredDate time.Time       `json:"registered_date" bson:"registered_date"`
	ActivePackage  Package         `json:"package" bson:"package"`
	Courses        []CoursesRecord `json:"courses_record" bson:"courses_record"`
	Country        string          `json:"country" bson:"country"`
	LastLogin      time.Time       `json:"last_login" bson:"last_login"`
	Parent         Parent          `json:"parent" bson:"parent"`
}

type StudentSummary struct {
	StudentId string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	Grade     uint8  `json:"grade" bson:"grade"`
	Country   string `json:"country" bson:"country"`
}

type StudentTrial struct {
	Name       string `json:"name" bson:"name"`
	Age        uint8  `json:"age" bson:"age"`
	Country    string `json:"country" bson:"country"`
	Experience uint8  `json:"exp" bson:"exp"`
}

type Parent struct {
	Name        string `json:"name" bson:"name"`
	PhoneNumber string `json:"phone" bson:"phone"`
	Email       string `json:"email" bson:"email"`
}

type CoursesRecord struct {
	CourseName string    `json:"course_name" bson:"course_name"`
	EnrolledAt time.Time `json:"enrolled_at" bson:"enrolled_at"`
	FinishedAt time.Time `json:"finished_at" bson:"finished_at"`
	Status     string    `json:"status" bson:"status"`
	Progress   float64   `json:"progress" bson:"progress"`
	Score      float64   `json:"score" bson:"score"`
}

type Package struct {
	CourseName string    `json:"course_name" bson:"course_name"`
	Status     string    `json:"status" bson:"status"`
	EnrolledAt time.Time `json:"enrolled_at" bson:"enrolled_at"`
}
