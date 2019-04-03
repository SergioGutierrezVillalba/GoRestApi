package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Timer struct{
	Id 		 	bson.ObjectId `bson:"_id"      json:"id,omitempty"`
	UserId	 	string  	  `bson:"userId"   json:"userId"`
	Start	 	int64         `bson:"start"	   json:"start"`
	Finish 	 	int64         `bson:"finish"   json:"finish"`
	Duration    int64	      `bson:"duration" json:"duration"`
}