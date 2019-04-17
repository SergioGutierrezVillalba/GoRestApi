package entity

import (
	"FirstProject/Model"
)

type RepositoryInterface interface {

	GetTasksOnTheSameDateAsUserTimers(string)([]model.TaskRepo, error)
	
}