package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Timer struct{
	Id 		 	bson.ObjectId `bson:"_id" json:"id,omitempty"`
	UserId	 	bson.ObjectId `bson:"userId" json:"userId"`
	Start	 	int64         `bson:"start" json:"start"`
	Finish 	 	int64         `bson:"finish" json:"finish"`
	Duration    int64	      `bson:"duration" json:"duration"`
}

func (t *Timer) GetId() (id string){
	return t.Id.Hex()
}

func (t *Timer) IsAlreadyFinished() bool {
	if t.Finish > 0 {
		return true
	}
	return false
}

func (t *Timer) NotExists() bool {
	if t == nil {
		return true
	}
	return false
}