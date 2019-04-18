package usecase

import (
	repo "FirstProject/Domains/task/entity"
	"FirstProject/Model"
)

type Usecase interface {
	GetTasksOnTheSameDateAsUserTimers(userId string) ([]model.TaskRepo, error)
	GetTasksAfterDateGiven(date int64)([]model.TaskRepo, error)
}

type TasksUsecase struct {
	repo		repo.RepositoryInterface
}

func NewUsecase(r repo.RepositoryInterface) Usecase {
	return &TasksUsecase{
		repo:r,
	}
}

func (t *TasksUsecase) GetTasksOnTheSameDateAsUserTimers(userId string) ([]model.TaskRepo, error) {
	return t.repo.GetTasksOnTheSameDateAsUserTimers(userId)
}

func (t *TasksUsecase) GetTasksAfterDateGiven(date int64)([]model.TaskRepo, error){
	return t.repo.GetTasksAfterDateGiven(date)
}