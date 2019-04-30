package user

import (
	mgo "gopkg.in/mgo.v2"
	"log"

	"FirstProject/Model"
)

type UserCollection struct {
	Db 		*mgo.Database
}

func NewUserCollection(d *mgo.Database) UserCollection {
	return UserCollection {
		Db: d,
	}
}

func (u *UserCollection) IsEmpty() (isEmpty bool) {

	isEmpty = true
	noRegisters, err := u.Db.C("users").Count()

	if err != nil {
		log.Print("error counting no. of documents")
	}

	if noRegisters > 0 {
		isEmpty = false
	}
	return
}

func (u *UserCollection) InsertUser(user model.User) error {
	return u.Db.C("users").Insert(&user)
}