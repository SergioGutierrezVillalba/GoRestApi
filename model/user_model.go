package model

import (
	"FirstProject/entities"
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	
)

type UserModel struct {
	Db *mgo.Database
}

func (userModel UserModel) FindAll() ([]entities.User, error) {

	var users []entities.User

	err := userModel.Db.C("users").Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	} else {
		return users, nil
	}

}

func (userModel UserModel) Find(id string) (entities.User, error) {

	var user entities.User

	err := userModel.Db.C("users").FindId(bson.ObjectIdHex(id)).One(&user)

	return user, err
}

func (userModel UserModel) FindByUsername(username string) (entities.User, error) {

	var user entities.User

	err := userModel.Db.C("users").Find(bson.M{"username": username}).One(&user)

	return user, err
}

func (userModel UserModel) Create(user *entities.User) error {

	pwd := []byte(user.Password)
	user.Password = string(GeneratePasswordHashed(pwd))

	return userModel.Db.C("users").Insert(&user)
}

func (userModel UserModel) Delete(id string) error { 
	return userModel.Db.C("users").RemoveId(bson.ObjectIdHex(id))
}

func (userModel UserModel) Update(user *entities.User) error {
	return userModel.Db.C("users").UpdateId(user.Id, user)
}

func (userModel UserModel) Login(user *entities.User) (response bool, err error) {

	fmt.Println(user.Password)
	pwd := []byte(user.Password)

	response = false
	
	if userDb, err := userModel.FindByUsername(user.Username); err != nil {
		return response, err

	} else {
		
		err2 := bcrypt.CompareHashAndPassword([]byte(userDb.Password), pwd) // if err = nil, correct password

		if err2 == nil {
			fmt.Println("Logged")
			response = true
			return response, nil
		} else {
			fmt.Println("Wrong password")
			return response, nil
		}
	}
	// Debugging purposes
	// GeneratePasswordHashed(pwd)
}

func (userModel UserModel) GetEmailUser(username string) (entities.User, error){
	
	var user entities.User
	
	err:= userModel.Db.C("users").Find(bson.M{"username": username}).One(&user)

	return user, err
}

// Functionalities
func GeneratePasswordHashed(pwd []byte) []byte {

	// Hashing the password with the default cost of 10
	hashedPwd, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(hashedPwd))
	return hashedPwd

}
