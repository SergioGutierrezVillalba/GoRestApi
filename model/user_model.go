package model 

import(
	"FirstProject/entities"
	
	"gopkg.in/mgo.v2/bson"
	mgo "gopkg.in/mgo.v2" 
)

type UserModel struct {
	Db *mgo.Database
}

func (userModel UserModel) FindAll() ([] entities.User, error){

	var users []entities.User

	err:= userModel.Db.C("users").Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	} else {
		return users, nil
	}
	
}

func (userModel UserModel) Find(id string) (entities.User, error){
	var user entities.User

	err:= userModel.Db.C("users").FindId(bson.ObjectIdHex(id)).One(&user)

	return user, err
}

func (userModel UserModel) Create(user *entities.User) error {
	return userModel.Db.C("users").Insert(&user)
}

func (userModel UserModel) Delete(id string) error {
	return userModel.Db.C("users").RemoveId(bson.ObjectIdHex(id))
}

func (userModel UserModel) Update(user *entities.User) error {
	return userModel.Db.C("users").UpdateId(user.Id, user)
}
