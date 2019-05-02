package entity

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"FirstProject/Model"
	"log"
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

func (m *MongoDbRepository) GetAll() (timers []model.Timer, err error){
	err = m.Db.C("times").Find(bson.M{}).All(&timers)
	return
}

func (m *MongoDbRepository) GetById(timerId string) (timer model.Timer, err error){
	err = m.Db.C("times").FindId(bson.ObjectIdHex(timerId)).One(&timer)
	return 
}

func (m *MongoDbRepository) GetAllByUserId(userId string) (timers []model.Timer, err error){
	err = m.Db.C("times").Find(bson.M{"userId": bson.ObjectIdHex(userId)}).All(&timers)
	return
}

func (m *MongoDbRepository) Create(timer model.Timer) error{
	log.Print(timer)
	return m.Db.C("times").Insert(&timer)
}

func (m *MongoDbRepository) Update(timer model.Timer) error{
	return m.Db.C("times").UpdateId(timer.Id, timer)
}

func (m *MongoDbRepository) Delete(timerId string) error{
	return m.Db.C("times").RemoveId(bson.ObjectIdHex(timerId))
}

func (m *MongoDbRepository) InsertStartTime(timer model.Timer) error{
	return m.Db.C("times").Insert(&timer)
}