package entity

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"FirstProject/Model"
	"fmt"
)

type MongoDbRepository struct{
	Session *mgo.Session
	Db 		*mgo.Database
}

func NewMongoDbRepository(s *mgo.Session) RepositoryInterface {
	return &MongoDbRepository{
		Session: s,
		Db: s.DB("project"),
	}
}

func (m *MongoDbRepository) GetAll() (users []model.User, err error){
	err = m.Db.C("users").Find(bson.M{}).All(&users)
	return
}

func (m *MongoDbRepository) GetById(userId string) (user model.User, err error){
	err = m.Db.C("users").FindId(bson.ObjectIdHex(userId)).One(&user)
	return 
}

func (m *MongoDbRepository) GetUserByJwt(jwt string)(user model.User, err error){
	err = m.Db.C("users").Find(bson.M{"jwt":jwt}).One(&user)
	return
}

func (m *MongoDbRepository) GetUserByUsername(username string)(user model.User, err error){
	err = m.Db.C("users").Find(bson.M{"username": username}).One(&user)
	return
}

func (m *MongoDbRepository) GetUserByRecoverToken(token string)(user model.User, err error){
	err = m.Db.C("users").Find(bson.M{"token": token}).One(&user)
	return
}

func (m *MongoDbRepository) Create(user model.User) error{
	return m.Db.C("users").Insert(&user)
}

func (m *MongoDbRepository) Update(user model.User) error{
	return m.Db.C("users").UpdateId(user.Id, user)
}

func (m *MongoDbRepository) Delete(userId string) error{
	return m.Db.C("users").RemoveId(bson.ObjectIdHex(userId))
}

func (m *MongoDbRepository) GetTasksOnTheSameDateAsUserTimers(userId string)([]model.Task, error) {

	var tasks []model.Task
	pipe := m.Db.C("tasks").Pipe([]bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":			"times",
				"localField": 	"timerId",
				"foreignField":	"_id",
				"as":			"timers",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":			"users",
				"localField": 	"timers.userId",
				"foreignField":	"_id",
				"as":			"userTasks",
			},
		},
		// bson.M{
		// 	"$project": bson.M{
		// 		"timers": 0,
		// 	},
		// },
		// bson.M{
		// 	"$match": bson.M{ 
		// 		"_id":bson.ObjectIdHex(userId),
		// 	},
		// },
	})

	err := pipe.All(&tasks)
	fmt.Println(tasks)
	return tasks, err
}