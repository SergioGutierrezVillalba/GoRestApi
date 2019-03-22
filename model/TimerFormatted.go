package model

import (
	"gopkg.in/mgo.v2/bson"
)

type TimerFormatted struct{
	Id 		 	bson.ObjectId `bson:"_id"      json:"id,omitempty"`
	UserId	 	string  	  `bson:"userId"   json:"userId"`
	Start	 	string        `bson:"start"	   json:"start"`
	Finish 	 	string        `bson:"finish"   json:"finish"`
	Duration    int64	      `bson:"duration" json:"duration"`
}