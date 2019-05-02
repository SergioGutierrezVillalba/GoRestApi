package controller

import (
	tasksUsecase "FirstProject/Domains/task/usecase"
	usersUsecase "FirstProject/Domains/user/usecase"

	"FirstProject/Model"
	"FirstProject/Model/Auth"
	"FirstProject/Model/helper"
	"FirstProject/Model/Sorter"
	
	"fmt"
	"time"
	"net/http"
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

type Controller interface{

	GetTasks(w http.ResponseWriter, r *http.Request)
	GetTaskById(w http.ResponseWriter, r *http.Request)
	GetTasksByTimerId(w http.ResponseWriter, r *http.Request)
	GetTasksSortedByCreationDate(w http.ResponseWriter, r *http.Request)
	GetTasksDoneByUserId(w http.ResponseWriter, r *http.Request)
	GetTasksDoneByUserIdSortedDescendent(w http.ResponseWriter, r *http.Request)
	GetNumberOfTasksAfterDateGiven(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)

}

var (
	Sorter			sorter.Sorter

	respond 		model.Responser

	Helper			helper.Helper

	jwtSent			string

	authenticator 	auth.Authentication
)

type TasksController struct {
	TasksUsecase		tasksUsecase.Usecase
	UsersUsecase		usersUsecase.Usecase
}

func NewController(t tasksUsecase.Usecase) Controller {
	return &TasksController{
		TasksUsecase: t,
	}
}

func (t *TasksController) GetTasks(w http.ResponseWriter, r *http.Request){

	tasks, err := t.TasksUsecase.GetTasks()

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, tasks)
}

func (t *TasksController) GetTaskById(w http.ResponseWriter, r *http.Request){

	taskId := GetIdFromUrl(r)
	task, err := t.TasksUsecase.GetTaskById(taskId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, task)
}

func (t *TasksController) GetTasksByTimerId(w http.ResponseWriter, r *http.Request){
	
	timerId := GetIdFromUrl(r)
	tasks, err := t.TasksUsecase.GetTasksByTimerId(timerId)

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, "GetTasksError")
		return
	}

	respond.WithJson(w, http.StatusOK, tasks)
}

func (t *TasksController) GetTasksSortedByCreationDate(w http.ResponseWriter, r *http.Request) {
	
	tasks, err := t.TasksUsecase.GetTasks()

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	tasksSliceSorted := Sorter.SortTasksSliceByCreationDate(tasks)
	respond.WithJson(w, http.StatusOK, tasksSliceSorted)
}

func (t *TasksController) GetTasksDoneByUserId(w http.ResponseWriter, r *http.Request){

	userId := GetIdFromUrl(r)
	tasks, err := t.TasksUsecase.GetTasksDoneByUserId(userId)

	fmt.Println(tasks)

	if ActionGivesError(err){
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	tasksSorted := Sorter.SortTasksSliceByCreationDate(tasks)
	respond.WithJson(w, http.StatusOK, tasksSorted)
}

func (t *TasksController) GetTasksDoneByUserIdSortedDescendent(w http.ResponseWriter, r *http.Request){

	userId := GetIdFromUrl(r)
	tasks, err := t.TasksUsecase.GetTasksDoneByUserId(userId)

	fmt.Println(tasks)

	if ActionGivesError(err){
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	tasksSorted := Sorter.SortTasksSliceByCreationDateDescendent(tasks)
	respond.WithJson(w, http.StatusOK, tasksSorted)
}

func (t *TasksController) GetNumberOfTasksAfterDateGiven(w http.ResponseWriter, r *http.Request) {

	date := GetDateFromUrl(r)
	tasks, err := t.TasksUsecase.GetNumberOfTasksAfterDateGiven(date)

	if ActionGivesError(err){
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, tasks)
}

func (t *TasksController) CreateTask(w http.ResponseWriter, r *http.Request){

	var task model.Task
	task.Id = bson.NewObjectId()
	task.CreationDate = time.Now().Unix()

	GetDataFromBodyJSONRequest(r, &task)

	if ActionGivesError(t.TasksUsecase.CreateTask(task)) {
		respond.WithError(w, http.StatusBadRequest, "CreateTaskError")
		return
	}

	respond.WithJson(w, http.StatusOK, task)
}

func (t *TasksController) UpdateTask(w http.ResponseWriter, r *http.Request){

	var task model.Task
	GetDataFromBodyJSONRequest(r, &task)

	taskUpdated, err := t.TasksUsecase.UpdateTask(task)
	
	if ActionGivesError(err){
		respond.WithError(w, http.StatusBadRequest, "UpdateTaskError")
		return
	}
	respond.WithJson(w, http.StatusOK, taskUpdated)
}

func (t *TasksController) DeleteTask(w http.ResponseWriter, r *http.Request){

	taskId := GetIdFromUrl(r)
	err := t.TasksUsecase.DeleteTask(taskId)

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}


// Func
func GetIdFromUrl(r *http.Request) (id string) {
	id = Helper.GetIdFromUrl(r)
	return
}

func GetDateFromUrl(r *http.Request) (date int64) {
	date = Helper.GetDateFromUrl(r)
	return
}

func GetDataFromBodyJSONRequest(r *http.Request, dataSaver interface{}){
	json.NewDecoder(r.Body).Decode(dataSaver)
}

func GetDataFromHeaderRequest(r *http.Request){
	jwtSent, _ = Helper.GetJWTFromHeaderRequest(r)
}

func ActionGivesError(e error) bool {
	return Helper.ActionGivesError(e)
}