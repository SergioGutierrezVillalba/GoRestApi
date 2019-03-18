package entities 

import (
	// "fmt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"FirstProject/model"
	"FirstProject/model/database"
)

type Repository struct {
	Db *mgo.Database
}

func (repository *Repository) GetUsers() ([]model.User, error){

	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return nil, err
	}

	var users []model.User

	if err := db.C("users").Find(bson.M{}).All(&users); err != nil {
		return nil, err
	} 

	return users, nil

}

func (repository *Repository) GetUserById(userId string) (interface{}, error){

	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return nil, err
	}

	var user model.User
	
	err = db.C("users").FindId(bson.ObjectIdHex(userId)).One(&user);

	return user, err
}

func (repository *Repository) GetUserByUsername(username string) (model.User, error) {

	db, err := GetSession()
	defer CloseSession(db)

	var user model.User

	if err != nil {
		return user, err
	}

	err = db.C("users").Find(bson.M{"username": username}).One(&user)

	return user, err
}

func (repository *Repository) GetEmailOfUser(username string) (model.User, error){

	db, err := GetSession()
	defer CloseSession(db)

	var user model.User

	if err != nil {
		return user, err
	}
	
	err = db.C("users").Find(bson.M{"username": username}).One(&user)

	return user, err
}

func (repository *Repository) GetUserByToken(token string) (model.User, error){

	db, err := GetSession()
	defer CloseSession(db)

	var user model.User

	if err != nil {
		return user, err
	}
	
	err = db.C("users").Find(bson.M{"token": token}).One(&user)

	return user, err
}

func (repository *Repository) GetUserByJwt(jwt string) (model.User, error){

	db, err := GetSession()
	defer CloseSession(db)

	var user model.User

	if err != nil {
		return user, err
	}
	
	err = db.C("users").Find(bson.M{"jwt": jwt}).One(&user)

	return user, err
}

func (repository *Repository) FindJwt(jwt string) error {

	db, err := GetSession()
	defer CloseSession(db)

	var user model.User

	if err != nil {
		return err
	}

	err2 := db.C("users").Find(bson.M{"jwt": jwt}).One(&user)
	return err2
}

func (repository *Repository) CreateUser(user model.User) error{

	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return err
	}

	return db.C("users").Insert(&user)
}

func (repository *Repository) UpdateUser(user model.User) error{

	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return err
	}

	return db.C("users").UpdateId(user.Id, user)
}

func (repository *Repository) DeleteUser(userId string) error{

	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return err
	}

	return db.C("users").RemoveId(bson.ObjectIdHex(userId))
}


// Functionalities
func GetSession()(*mgo.Database, error){
	return database.GetSession()
}

func CloseSession(db *mgo.Database){
	database.CloseSession(db)
}
