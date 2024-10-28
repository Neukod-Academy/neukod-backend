package models

import "time"

type Trial struct {
	TrialId    string       `json:"trial_id" bson:"trial_id"`
	ParentName string       `json:"parent_name" bson:"parent_name"`
	Student    StudentTrial `json:"student" bson:"student"`
	Note       string       `json:"note" bson:"note"`
	CreatedAt  time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" bson:"updated_at"`
}
