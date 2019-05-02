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

func (m *MongoDbRepository) GetTasks()(tasks []model.Task, err error) {
	err = m.Db.C("tasks").Find(bson.M{}).All(&tasks)
	return
}

func (m *MongoDbRepository) GetTaskById(taskId string)(task model.Task, err error) {
	err = m.Db.C("tasks").FindId(bson.ObjectIdHex(taskId)).One(&task)
	return 	
}

func (m *MongoDbRepository) GetTasksByTimerId(timerId string)(tasks []model.Task, err error) {
	err = m.Db.C("tasks").Find(bson.M{"timerId": bson.ObjectIdHex(timerId)}).All(&tasks)
	return
}

func (m *MongoDbRepository) GetTasksDoneByUserId(userId string)([]model.Task, error) {

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

	err := pipe.All(&tasks)
	return tasks, err
}

func (m *MongoDbRepository) GetNumberOfTasksAfterDateGiven(date int64)([]model.TaskRepo, error){

	var tasksRepo []model.TaskRepo
	
	pipe := m.Db.C("tasks").Pipe([]bson.M{
		bson.M{
			"$match": bson.M{ 
				"$expr": bson.M { 
					"$gte": []interface{}{"$creationDate", date},
				},
			},
		},
		bson.M{
			"$count": "tasks_done_after_date",
		},
		bson.M{
			"$project": bson.M{
				"timers":0,
			},
		},
	})

	err := pipe.All(&tasksRepo)
	return tasksRepo, err
}

func (m *MongoDbRepository) CreateTask(task model.Task) error {
	return m.Db.C("tasks").Insert(&task)
}

func (m *MongoDbRepository) UpdateTask(task model.Task) error {	
	return m.Db.C("tasks").UpdateId(task.Id, task)
}

func (m *MongoDbRepository) DeleteTask(taskId string) error {
	return m.Db.C("tasks").RemoveId(bson.ObjectIdHex(taskId))	
}