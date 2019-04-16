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