package models

import (
	"fmt"
	"time"
)

var TrialList = []string{
	"roblox",
	"c++",
	"c#",
	"java",
}

type Trial struct {
	TrialId   string       `json:"trial_id" bson:"trial_id"`
	Parent    Parent       `json:"parent" bson:"parent"`
	Student   StudentTrial `json:"student" bson:"student"`
	Course    string       `json:"course" bson:"course"`
	Note      string       `json:"note" bson:"note"`
	CreatedAt time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" bson:"updated_at"`
}

func (t *Trial) CheckIfEmpty() []string {
	var errors []string
	if t.Parent.Name == "" {
		errors = append(errors, "parent name is required")
	}
	if t.Parent.Email == "" {
		errors = append(errors, "parent email is required")
	}
	if t.Parent.PhoneNumber == "" {
		errors = append(errors, "parent phone number is required")
	}
	if t.Student.Name == "" {
		errors = append(errors, "student name is required")
	}
	if t.Student.Age == 0 {
		errors = append(errors, "student age is required")
	}
	if t.Student.Experience == 0 {
		errors = append(errors, "student experience is required")
	}
	if t.Course == "" {
		errors = append(errors, "course is required")
	}
	fmt.Println(errors)
	return errors
}
