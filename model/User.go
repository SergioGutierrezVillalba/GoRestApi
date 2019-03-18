package model

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct{
	Id 		 	bson.ObjectId `bson:"_id"      json:"id"`
	Username 	string  	  `bson:"username" json:"username"`
	Password 	string        `bson:"password" json:"password"`
	Role 	 	string        `bson:"role"     json:"role"`
	Email    	string	      `bson:"email"    json:"email"`
	Token    	string	      `bson:"token"    json:"token"`
	Jwt      	string		  `bson:"jwt"      json:"jwt"`
}