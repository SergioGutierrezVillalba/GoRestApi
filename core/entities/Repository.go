package entities 

import (
	
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"errors"

	"FirstProject/model"
	"FirstProject/model/database"
)

type Repository struct {
	Db *mgo.Database
}

// USERS

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
		return user, errors.New("DatabaseError")
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
		return errors.New("UpdateUserError")
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


// TIMER

func (repository *Repository) GetTimers() ([] model.Timer, error){

	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return nil, err
	}

	var timers []model.Timer

	if err := db.C("times").Find(bson.M{}).All(&timers); err != nil {
		return nil, err
	} 

	return timers, nil
}

func (repository *Repository) GetTimeById(timeId string) (model.Timer, error) {

	var timer model.Timer
	
	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return timer, err
	}

	err = db.C("times").FindId(bson.ObjectIdHex(timeId)).One(&timer);

	return timer, err
}

func (repository *Repository) InsertStartTime(time model.Timer) error {

	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return err
	}

	if err := db.C("times").Insert(&time); err != nil {
		return err
	}

	return nil
}

func (repository *Repository) InsertFinishTime(time model.Timer) error {

	db, err := GetSession()
	defer CloseSession(db)

	if err != nil {
		return errors.New("UpdateTimeError")
	}

	return db.C("times").UpdateId(time.Id, time)
}



// Functionalities
func GetSession()(*mgo.Database, error){
	return database.GetSession()
}

func CloseSession(db *mgo.Database){
	database.CloseSession(db)
}
