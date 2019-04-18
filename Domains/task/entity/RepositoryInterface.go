package entity

import (
	"FirstProject/Model"
)

type RepositoryInterface interface {
	GetTasksOnTheSameDateAsUserTimers(string)([]model.TaskRepo, error)
	GetTasksAfterDateGiven(date int64)([]model.TaskRepo, error)
}