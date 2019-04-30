package usecase

import (
	"errors"

	repo "FirstProject/Domains/task/entity"
	"FirstProject/Model"
	"FirstProject/Model/Validation"
)

type Usecase interface {
	GetTasks()([]model.Task, error)
	GetTaskById(string)(model.Task, error)
	GetTasksByTimerId(string)([]model.Task, error)
	GetTasksDoneByUserId(string)([]model.Task, error)
	GetNumberOfTasksAfterDateGiven(date int64)([]model.TaskRepo, error)
	CreateTask(model.Task) error
	UpdateTask(model.Task) (model.Task, error)
	DeleteTask(string) error
}

var checker validation.Checker

type TasksUsecase struct {
	repo		repo.RepositoryInterface
}

func NewUsecase(r repo.RepositoryInterface) Usecase {
	return &TasksUsecase{
		repo:r,
	}
}

func (t *TasksUsecase) GetTasks()(tasks []model.Task, err error){
	return t.repo.GetTasks()
}

func (t *TasksUsecase) GetTaskById(taskId string)(model.Task, error){
	return t.repo.GetTaskById(taskId)
}

func (t *TasksUsecase) GetTasksByTimerId(timerId string)([]model.Task, error){
	return t.repo.GetTasksByTimerId(timerId)
}

func (t *TasksUsecase) GetTasksDoneByUserId(userId string) ([]model.Task, error) {
	return t.repo.GetTasksDoneByUserId(userId)
}

func (t *TasksUsecase) GetNumberOfTasksAfterDateGiven(date int64)([]model.TaskRepo, error){
	return t.repo.GetNumberOfTasksAfterDateGiven(date)
}

func (t *TasksUsecase) CreateTask(task model.Task) error {
	return t.repo.CreateTask(task)
}

func (t *TasksUsecase) UpdateTask(taskUpdated model.Task) (model.Task, error) {

	taskDb, err := t.GetTaskById(taskUpdated.Id.Hex())

	if err != nil {
		return taskUpdated, errors.New("TaskNotExistError")
	}

	taskUpdated.CreationDate = taskDb.CreationDate
	return taskUpdated, t.repo.UpdateTask(taskUpdated)
}

func (t *TasksUsecase) DeleteTask(taskId string) error {
	return t.repo.DeleteTask(taskId)
}