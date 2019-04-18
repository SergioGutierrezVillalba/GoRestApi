package controller

import (
	tasksUsecase "FirstProject/Domains/task/usecase"

	"FirstProject/Model"
	"FirstProject/Model/helper"

	"net/http"
)

type Controller interface{

	GetTasksOnTheSameDateAsUserTimersByUserId(w http.ResponseWriter, r *http.Request)
	GetTasksAfterDateGiven(w http.ResponseWriter, r *http.Request)

}

var (
	respond 		model.Responser
	Helper			helper.Helper
)

type TasksController struct {
	TasksUsecase		tasksUsecase.Usecase
}

func NewController(t tasksUsecase.Usecase) Controller {
	return &TasksController{
		TasksUsecase: t,
	}
}

func (t *TasksController) GetTasksOnTheSameDateAsUserTimersByUserId(w http.ResponseWriter, r *http.Request){
	userId := GetIdFromUrl(r)
	tasks, err := t.TasksUsecase.GetTasksOnTheSameDateAsUserTimers(userId)

	if ActionGivesError(err){
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.WithJson(w, http.StatusOK, tasks)
}

func (t *TasksController) GetTasksAfterDateGiven(w http.ResponseWriter, r *http.Request) {
	var date int64
	date = 1555338589 // lunes, 15 de abril de 2019 16:29:49 GMT+02:00
	tasks, err := t.TasksUsecase.GetTasksAfterDateGiven(date)

	if ActionGivesError(err){
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, tasks)
}

// Func
func GetIdFromUrl(r *http.Request) (id string) {
	id = Helper.GetIdFromUrl(r)
	return
}

func ActionGivesError(e error) bool {
	return Helper.ActionGivesError(e)
}