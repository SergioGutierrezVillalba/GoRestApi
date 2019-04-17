package entity

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"FirstProject/Model"
)

type MongoDbRepository struct {
	Session *mgo.Session
	Db 		*mgo.Database
}

func NewMongoDbRepository(s *mgo.Session) RepositoryInterface {
	return &MongoDbRepository{
		Session: s,
		Db: s.DB("project"),
	}
}

func (m *MongoDbRepository) GetTasksOnTheSameDateAsUserTimers(userId string)([]model.TaskRepo, error) {

	var tasksRepo []model.TaskRepo
	
	pipe := m.Db.C("tasks").Pipe([]bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":			"times",
				"localField": 	"timerId",
				"foreignField":	"_id",
				"as":			"timers",
			},
		},
		{
			"$unwind" : "$timers",
		},
		bson.M{
			"$lookup": bson.M{
				"from":			"users",
				"localField": 	"timers.userId",
				"foreignField":	"_id",
				"as":			"timers.users",
			},
		},
		bson.M{
			"$project": bson.M{
				"timers.users": 0,
			},
		},
		bson.M{
			"$match": bson.M{ 
				"timers.userId": bson.ObjectIdHex(userId),
				"$expr": bson.M { "$and": []bson.M{
						{"$gte" : []string{"$creationDate","$timers.start"}},
						{"$lte" : []string{"$creationDate","$timers.finish"}},
					},
				},
			},
		},
	})

	err := pipe.All(&tasksRepo)
	return tasksRepo, err
}