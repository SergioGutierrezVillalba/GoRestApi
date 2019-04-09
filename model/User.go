package model

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct{
	Id 		 	bson.ObjectId `bson:"_id"      json:"id,omitempty"`
	Username 	string  	  `bson:"username" json:"username"`
	Password 	string        `bson:"password" json:"password"`
	Role 	 	string        `bson:"role"     json:"role"`
	Email    	string	      `bson:"email"    json:"email"`
	Token    	string	      `bson:"token"    json:"token"`
	Jwt      	string		  `bson:"jwt"      json:"jwt"`
	RawId		string		  `bson:"rawId"	   json:"raw"`
	GroupId		string		  `bson:"groupId"  json:"groupId"`

}

func (u *User) HasGroup() bool {
	if u.GroupId == "" {
		return false
	}
	return true
}

func (u *User) IsFromTheSameGroup(groupId string) bool {
	if u.GroupId == groupId {
		return true
	}
	return false
}