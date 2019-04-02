package controller

import (
	usersUsecase "FirstProject/Domains/user/usecase"
	timersUsecase "FirstProject/Domains/timer/usecase"
	"FirstProject/model/auth"
	"FirstProject/model"
	"FirstProject/model/helper"

	"net/http"
	"time"
	"strconv"
	"fmt"
	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

type Controller interface{
	
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	GetUserByJwt(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	SendRecover(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)

	GetAllTimers(w http.ResponseWriter, r *http.Request)
	GetTimerById(w http.ResponseWriter, r *http.Request)
	GetTimersByUserId(w http.ResponseWriter, r *http.Request)
	CreateTimer(w http.ResponseWriter, r *http.Request)
	UpdateTimer(w http.ResponseWriter, r *http.Request)
	DeleteTimer(w http.ResponseWriter, r *http.Request)
	StartTimer(w http.ResponseWriter, r *http.Request)
	FinishTimer(w http.ResponseWriter, r *http.Request)
}

type UsersController struct {
	UsersUsecase	 usersUsecase.Usecase
	TimersUsecase	 timersUsecase.Usecase
}

var (
	respond 		model.Responser
	mailSender		model.MailSender

	Helper			helper.Helper

	jwtSent			string

	crypter 		auth.Crypter
	authenticator 	auth.Authentication
	responseToken	auth.ResponseToken
	requestInfo 	auth.RequestInfo
)

func NewUsersController(u usersUsecase.Usecase, t timersUsecase.Usecase) Controller {
	return &UsersController{
		UsersUsecase: u,
		TimersUsecase: t,
	}
}

// USERS CONTEXT

func (u *UsersController) GetAllUsers(w http.ResponseWriter, r *http.Request){

	users, err := u.UsersUsecase.GetAll()

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	respond.WithJson(w, http.StatusOK, users)
}

func (u *UsersController) GetUserByJwt(w http.ResponseWriter, r *http.Request){
	
	GetDataFromHeaderRequest(r)
	user, err := u.UsersUsecase.GetUserByJwt(jwtSent);
	user.Password = ""
	user.Jwt = ""

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, user)
}

func (u *UsersController) GetUserById(w http.ResponseWriter, r *http.Request){

	userId := GetIdFromUrl(r)
	user, err := u.UsersUsecase.GetById(userId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, user)
}

func (u *UsersController) CreateUser(w http.ResponseWriter, r *http.Request){

	var user model.User
	user.Id = bson.NewObjectId()

	GetDataFromBodyRequest(r, &user)

	user.Password, _ = crypter.Crypt(user.Password)

	if err := u.UsersUsecase.Create(user); err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user.Password = ""
	respond.WithJson(w, http.StatusOK, user)
}

func (u *UsersController) UpdateUser(w http.ResponseWriter, r *http.Request){

	var userToUpdate model.User
	GetDataFromBodyRequest(r, &userToUpdate)
	GetDataFromHeaderRequest(r)

	// SELF VERIFY
	userIdRequesting, _ := authenticator.GetUserIdFromJWT(jwtSent) 
	userRequesting, _ := u.UsersUsecase.GetById(userIdRequesting)
	fmt.Println("(UpdateUser): Id user requesting: " + userIdRequesting)
	fmt.Println("(UpdateUser): Id user wants upda: " + userToUpdate.Id.Hex())

	// WHICH ROLE IS USING
	roleUsed := WhichRoleIsUsed(userRequesting, userToUpdate)
	switch roleUsed {
		case "NOAUTH":
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		case "SELF":
			userToUpdate.Role = "user"
		case "ADMIN":
			// make sth that you need if is an admin
	}

	savePWD := userToUpdate.Password
	userToUpdate.Password = ""

	newJWT := authenticator.GenerateJWT(userToUpdate)
	userToUpdate.Jwt = newJWT
	userToUpdate.Password, _ = crypter.Crypt(savePWD)

	err := u.UsersUsecase.Update(userToUpdate)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.WithJson(w, http.StatusOK, auth.ResponseToken{Token:userToUpdate.Jwt})
}

func (u *UsersController) DeleteUser(w http.ResponseWriter, r *http.Request){

	idUser := GetIdFromUrl(r)
	err:= u.UsersUsecase.Delete(idUser);

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController) Login(w http.ResponseWriter, r *http.Request){

	var user model.User

	GetDataFromBodyRequest(r, &user)
	userDb, err := u.UsersUsecase.GetUserByUsername(user.Username)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = crypter.PasswordCoincides(userDb.Password, user.Password)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	savedPassword := userDb.Password

	userDb.Password = ""
	userDb.Jwt = ""

	token := authenticator.GenerateJWT(userDb)

	userDb.Jwt = token
	userDb.Password = savedPassword

	err = u.UsersUsecase.Update(userDb)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("This is the token that im sending")
	fmt.Println(token)
	respond.WithJson(w, http.StatusOK, auth.ResponseToken{Token:token})
}

func (u *UsersController) Register(w http.ResponseWriter, r *http.Request){

	var user model.User
	user.Id = bson.NewObjectId()

	GetDataFromBodyRequest(r, &user)

	user.Password, _ = crypter.Crypt(user.Password)
	user.Role = "user"

	if err := u.UsersUsecase.Register(user); err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController) SendRecover(w http.ResponseWriter, r *http.Request){
	
	var user model.User
	GetDataFromBodyRequest(r, &user)

	userDb, err := u.UsersUsecase.GetUserByUsername(user.Username)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err2 := crypter.Crypt(userDb.Email)

	if err2 != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = mailSender.Send(userDb.Email, token)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController) ResetPassword(w http.ResponseWriter, r *http.Request){

	var passwordRecover model.PasswordRecover
	GetDataFromBodyRequest(r, &passwordRecover)

	user, err := u.UsersUsecase.GetUserByRecoverToken(passwordRecover.Token)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user.Password, err = crypter.Crypt(passwordRecover.NewPassword)
	user.Token = ""

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = u.UsersUsecase.Update(user); err != nil{
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	respond.WithJson(w, http.StatusOK, "Success")
}

// TIMERS CONTEXT

func (u *UsersController)GetAllTimers(w http.ResponseWriter, r *http.Request){
	
	timers, err := u.TimersUsecase.GetAll()

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted, err := formatTimersForResponse(timers)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, timersFormatted)
}

func (u *UsersController)GetTimerById(w http.ResponseWriter, r *http.Request){

	timerId := GetIdFromUrl(r)
	timer, err := u.TimersUsecase.GetById(timerId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted, err := formatTimerForResponse(timer)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController)GetTimersByUserId(w http.ResponseWriter, r *http.Request){

	userId := GetIdFromUrl(r)
	fmt.Println("Controller, userId: " + userId)
	timers, err := u.TimersUsecase.GetAllByUserId(userId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted, err := formatTimersForResponse(timers)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, timersFormatted)
	
}

func (u *UsersController)CreateTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	timer.Id = bson.NewObjectId()

	GetDataFromBodyRequest(r, &timer)

	err := u.TimersUsecase.Create(timer)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted, err := formatTimerForResponse(timer)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController)UpdateTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	GetDataFromBodyRequest(r, &timer)

	err := u.TimersUsecase.Update(timer)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted, err := formatTimerForResponse(timer)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController)DeleteTimer(w http.ResponseWriter, r *http.Request){

	timerId := GetIdFromUrl(r)
	err := u.TimersUsecase.Delete(timerId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController)StartTimer(w http.ResponseWriter, r *http.Request){

	var user model.User
	var timer model.Timer

	// SELF VERIFY
	GetDataFromBodyRequest(r, &user)
	GetDataFromHeaderRequest(r)

	userIdRequesting,_ := authenticator.GetUserIdFromJWT(jwtSent)
	userRequesting, _ := u.UsersUsecase.GetById(userIdRequesting)

	fmt.Println("(StartTimer): Id user requesting: " + userIdRequesting)
	fmt.Println("(StartTimer): Id user wants to start timer: " + user.Id.Hex())

	// CHECK ROLE USED
	roleUsed := WhichRoleIsUsed(userRequesting, model.User{Id: bson.ObjectId(timer.UserId)})
	switch roleUsed {
		case "NOAUTH":
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		case "SELF":
			// make sth that you need if is a self
		case "ADMIN":
			// make sth that you need if is an admin
	}

	timer.Id = bson.NewObjectId()
	timer.UserId = user.Id.Hex()
	timer.Start = time.Now().Unix()

	err := u.TimersUsecase.StartTimer(timer)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted, _ := formatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController)FinishTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	// GET DATA 
	GetDataFromBodyRequest(r, &timer)
	GetDataFromHeaderRequest(r)
	
	timerDb, err := u.TimersUsecase.GetById(timer.Id.Hex())
	// TIMER NOT EXIST
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// TIMER IS ALREADY FINISHED
	if timerDb.Finish > 0 {
		respond.WithError(w, http.StatusBadRequest, "TimerAlreadyFinished")
		return
	}

	userIdRequesting, _ := authenticator.GetUserIdFromJWT(jwtSent)
	userRequesting, _ := u.UsersUsecase.GetById(userIdRequesting)
	fmt.Println("(FinishTimer): Id user requesting: " + userIdRequesting)
	fmt.Println("(FinishTimer): Id user wants stop: " + timerDb.UserId)

	// IS ADMIN, SELF, NONE
	roleUsed := WhichRoleIsUsed(userRequesting, model.User{Id: bson.ObjectId(timer.UserId)})
	switch roleUsed {
		case "NOAUTH":
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		case "SELF":
			// make sth that you need if is a self
		case "ADMIN":
			// make sth that you need if is an admin
	}
	
	finishTime := time.Now().Unix()
	duration := finishTime - timerDb.Start

	timerDb.Finish = finishTime
	timerDb.Duration = duration

	err = u.TimersUsecase.Update(timerDb)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted, err2 := formatTimerForResponse(timerDb)

	if err2 != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, timerFormatted)
}


// Func
func formatTimerForResponse(timer model.Timer) (timerFormatted model.TimerFormatted, err error){

	timerFormatted.Id = timer.Id
	timerFormatted.UserId = timer.UserId
	timerFormatted.Duration = timer.Duration
	timerFormatted.Start, err = unixDateToString(timer.Start)

	if err != nil {
		return
	}

	timerFormatted.Finish, err = unixDateToString(timer.Finish)

	if err != nil {
		return
	}

	return
}

func formatTimersForResponse(timers [] model.Timer) (timersFormatted [] model.TimerFormatted, err error){

	for i := 0; i < len(timers); i++ {
		timerFormatted, errFormat := formatTimerForResponse(timers[i])

		if errFormat != nil {
			err = errFormat
			return
		}
		
		timersFormatted = append(timersFormatted, timerFormatted)
	}
	return
}

func unixDateToString(miliseconds int64) (string, error) {

	miliToString := strconv.FormatInt(miliseconds, 10) // int64 to String
	i, err := strconv.ParseInt(miliToString, 10, 64) // String to int

    if err != nil {
        return "", err
	}
	tm := time.Unix(i, 0) // int to Unix timestamp
	tmToString := tm.String()

	return tmToString, nil
}

func GetIdFromUrl(request *http.Request) (id string) {
	vars := mux.Vars(request)
	id = vars["id"]
	return
}

func GetDataFromBodyRequest(r *http.Request, dataSaver interface{}){
	json.NewDecoder(r.Body).Decode(dataSaver)
}

func GetDataFromHeaderRequest(r *http.Request){
	jwtSent = Helper.GetJWTFromHeader(r)
}

// We have to determine what user is Requesting and
// which is susceptible to change to know which Role
// has the user Requesting and in consequence what it
// is allowed to do

func WhichRoleIsUsed (userRequesting model.User, userToModify model.User) (situation string) {
	if Helper.IsUser(userRequesting) {
		if !Helper.IsUpdatingItself(userRequesting.Id.Hex(), userToModify.Id.Hex()){
			fmt.Println("(SelfRoleRequest): You arent allowed to do that")
			return "NOAUTH"
		}
		return "SELF"
	}
	return "ADMIN"
}