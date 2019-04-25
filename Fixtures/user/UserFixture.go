package user

import (
	"log"
	"encoding/json"
	"io/ioutil"
	
	"FirstProject/Model"
	usercol "FirstProject/Collection/user"
)

type UserFixture struct {	
	UserCollection		usercol.UserCollection
}

func NewUserFixture(u usercol.UserCollection) UserFixture {
	return UserFixture {
		UserCollection: u,
	}
}

func (u *UserFixture) LoadFixture(){
	if u.UserCollection.IsEmpty() {
		log.Print("(UserCollection: I'm empty. Refilling...)")
		usersFixture := GetUsersFixtureFromJSON()

		for _, user := range usersFixture {
			err := u.UserCollection.InsertUser(user)
			log.Print(user)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func GetUsersFixtureFromJSON()(usersFixture []model.User){

	file, err := ioutil.ReadFile("./usersfixture.json")
	if err != nil {
		log.Print("Error during usersfixture reading")
	}

	err2 := json.Unmarshal([]byte(file), &usersFixture)
	if err2 != nil {
		log.Print(err2)
	}
	return
}