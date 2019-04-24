package entity

import (
	"FirstProject/Model"
)

type RepositoryInterface interface {
	GetTasksOnTheSameDateAsUserTimers(string)([]model.TaskRepo, error)
	GetNumberOfTasksAfterDateGiven(date int64)([]model.TaskRepo, error)
}