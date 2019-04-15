package model

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct{
	Id 		 		bson.ObjectId `bson:"_id"          json:"id,omitempty"`
	Username 		string  	  `bson:"username"     json:"username"`
	Password 		string        `bson:"password"     json:"password"`
	Role 	 		string        `bson:"role"         json:"role"`
	Email    		string	      `bson:"email"        json:"email"`
	Token    		string	      `bson:"token"        json:"token"`
	Jwt      		string		  `bson:"jwt"          json:"jwt"`
	RawId			string		  `bson:"rawId"	       json:"raw"`
	GroupId			string		  `bson:"groupId"  	   json:"groupId"`
	ProfileImage 	string		  `bson:"profileImage" json:"profileImage"`
	RouteImg		string		  `bson:"routeimg"	   json:"routeimg"`
}

// Getters
func (u *User) GetId() (id string){
	id = u.Id.Hex()
	return
}

// Miscelanea

// Checks if username is empty 
// so the user arrived null
// so user not exists
func (u *User) NotExists() bool {
	if u.Username == "" {
		return true
	}
	return false
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

// Setters
func (u *User) SetRouteImg(routeImg string){
	u.RouteImg = routeImg
}
func (u *User) SetJWT(jwt string){
	u.Jwt = jwt
}
func (u *User) SetPassword(newPassword string){
	u.Password = newPassword
}
func (u *User) SetRole(role string){
	u.Role = role
}

// Cleaners
func (u *User) EmptyPassword(){
	u.Password = ""
}
func (u *User) EmptyJWT(){
	u.Jwt = ""
}
func (u *User) EmptyProfileImage(){
	u.ProfileImage = ""
}