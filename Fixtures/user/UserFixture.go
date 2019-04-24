package user

import (
	"log"
	"encoding/json"
	"io/ioutil"
	
	"FirstProject/Model"
)

type UserFixture struct {	
}

func (u *UserFixture) LoadUsersFixture(){
	usersFixture := GetUsersFixtureFromJSON()
	for _, user := range usersFixture {
		log.Print(user)
	}
}

func GetUsersFixtureFromJSON()(usersFixture []model.User){

	file, err := ioutil.ReadFile("../../usersfixture.json")
	if err != nil {
		log.Print("Error during usersfixture reading")
	}

	err2 := json.Unmarshal([]byte(file), &usersFixture)
	if err2 != nil {
		log.Print("Error during usersfixture conversion")
	}
	return
}