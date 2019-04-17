package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Task struct {
	Id 		 		bson.ObjectId   `json:"id,omitempty" bson:"_id"`
	TimerId			bson.ObjectId	`json:"timerId" bson:"timerId"`
	Content			string			`json:"content"	bson:"content"`
	CreationDate	int64			`json:"creationDate" bson:"creationDate"`
}

type TaskRepo struct {
	Id 		 		bson.ObjectId   `json:"id,omitempty" bson:"_id"`
	TimerId			bson.ObjectId	`json:"timerId" bson:"timerId"`
	Content			string			`json:"content"	bson:"content"`
	CreationDate	int64			`json:"creationDate" bson:"creationDate"`
	TasksDone		int64			`json:"tasks_done_after_date" bson:"tasks_done_after_date"`
	Timers			TimerRepo		`json:"timers" bson:"timers,omitempty"`
}

type TimerRepo struct {
	Id 		 	bson.ObjectId `bson:"_id" json:"id,omitempty"`
	UserId	 	bson.ObjectId `bson:"userId" json:"userId"`
	Start	 	int64         `bson:"start" json:"start"`
	Finish 	 	int64         `bson:"finish" json:"finish"`
	Duration    int64	      `bson:"duration" json:"duration"`
	Users		[]User		  `json:"users,omitempty" bson:"users"`
}