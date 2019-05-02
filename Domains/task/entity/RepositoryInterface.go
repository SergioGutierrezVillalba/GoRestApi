package entity

import (
	"FirstProject/Model"
)

type RepositoryInterface interface {
	GetTasks()([]model.Task, error)
	GetTaskById(string)(model.Task, error)
	GetTasksByTimerId(string)([]model.Task, error)
	GetTasksDoneByUserId(string)([]model.Task, error)
	GetNumberOfTasksAfterDateGiven(date int64)([]model.TaskRepo, error)
	CreateTask(model.Task) error
	UpdateTask(model.Task) error
	DeleteTask(string) error
}