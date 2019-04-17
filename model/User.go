package model

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct{
	Id 		 		bson.ObjectId `bson:"_id" json:"id,omitempty"`
	Username 		string  	  `bson:"username" json:"username,omitempty"`
	Password 		string        `bson:"password" json:"password,omitempty"`
	Role 	 		string        `bson:"role" json:"role,omitempty"`
	Email    		string	      `bson:"email" json:"email,omitempty"`
	Token    		string	      `bson:"token" json:"token,omitempty"`
	Jwt      		string		  `bson:"jwt" json:"jwt,omitempty"`
	RawId			string		  `bson:"rawId" json:"raw,omitempty"`
	GroupId			string		  `bson:"groupId" json:"groupId,omitempty"`
	ProfileImage 	string		  `bson:"profileImage" json:"profileImage,omitempty"`
	RouteImg		string		  `bson:"routeimg" json:"routeimg,omitempty"`
}

// Getters
func (u *User) GetId() (id string){
	id = u.Id.Hex()
	return
}

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