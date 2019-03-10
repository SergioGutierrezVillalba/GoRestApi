package model

import (
	"log"
	mgo "gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
)

type Repository struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "users"
)

func(r *Repository) Connect() {

	session, err := mgo.Dial(r.Server)

	if err != nil {
		log.Fatal(err)
	}

	db = session.DB(r.Database)
}


// func (r *Repository) GetAll() ([] User, error) {
// 	var users [] User
// 	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
// 	return users, err
// }


// func (r *Repository) GetById(id string) (User, error) {
// 	var user User
// 	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
// 	return user, err
// }


// func (r *Repository) InsertUser(user User) error {
// 	err := db.C(COLLECTION).Insert(&user)
// 	return err
// }


// func (r *Repository) DeleteUser(user User) error {
// 	err := db.C(COLLECTION).Remove(&user)
// 	return err
// }


// func (r *Repository) UpdateUser(user User) error {
// 	err := db.C(COLLECTION).UpdateId(user.ID, &user)
// 	return err
// }