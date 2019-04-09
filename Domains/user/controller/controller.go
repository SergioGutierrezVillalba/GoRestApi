package controller

import (
	usersUsecase "FirstProject/Domains/user/usecase"
	timersUsecase "FirstProject/Domains/timer/usecase"

	"FirstProject/model/auth"
	"FirstProject/model/helper"
	"FirstProject/model/socket"
	"FirstProject/model/mail"
	"FirstProject/model"

	"net/http"
	"time"
	"fmt"
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
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
	
	StartWebSocket(w http.ResponseWriter, r *http.Request)
	FinishWebSocket(w http.ResponseWriter, r *http.Request)

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
	socketsMaps		= make(map[string]map[string]*socket.WebSocket)

	respond 		model.Responser
	
	mailSender		mail.MailSender

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

// SOCKETS CONTEXT

func(u *UsersController) StartWebSocket(w http.ResponseWriter, r *http.Request) {

	ws, err := socket.NewWebSocket(w, r)
	var tokenReceived string

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ws.On("message", func(e *socket.Event) {
		tokenReceived = e.Data.(string)
		uncryptedJWT := authenticator.Decrypt(authenticator.DecodeBase64(tokenReceived))
		userFromJWT := authenticator.GetUserInfoFromJWT(uncryptedJWT)
		fmt.Println("(StartWebSocket): GroupId - " + userFromJWT.GroupId)

		idMap := ksuid.New().String()
		groupAndSocket := make(map[string]*socket.WebSocket)
		groupAndSocket[userFromJWT.GroupId] = ws
		socketsMaps[idMap] = groupAndSocket

		// log.Printf("[MESSAGE] %v", e.Data)
		ws.Out <- (&socket.Event{
			Name: "response",
			Data: "Socket created |" + idMap,
		}).Raw()
	})
}

func (u *UsersController) FinishWebSocket(w http.ResponseWriter, r *http.Request){

	var socket socket.SocketResponser
	GetDataFromBodyRequest(r, &socket)
	// Get that websocket and close conn
	socketsMaps[socket.SocketId][socket.GroupId].Conn.Close()
	// Delete that websocketMap from generalMap
	delete(socketsMaps, socket.SocketId)

	// Will work with ws.on("close", {close conn and delete from websocketMap})?
	// Must be that code put in StartWebSocket method with 'on' event?
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
	roleUsed := Helper.WhichRoleIsUsed(userRequesting, userToUpdate)
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

	userDb.Token = token

	if err = u.UsersUsecase.Update(userDb); err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = mailSender.SendRecover(userDb.Email, token)

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

func (u *UsersController) GetAllTimers(w http.ResponseWriter, r *http.Request){
	
	timers, err := u.TimersUsecase.GetAll()

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted := Helper.FormatTimersForResponse(timers)
	respond.WithJson(w, http.StatusOK, timersFormatted)
}

func (u *UsersController) GetTimerById(w http.ResponseWriter, r *http.Request){

	timerId := GetIdFromUrl(r)
	timer, err := u.TimersUsecase.GetById(timerId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController) GetTimersByUserId(w http.ResponseWriter, r *http.Request){

	userId := GetIdFromUrl(r)
	fmt.Println("Controller, userId: " + userId)
	timers, err := u.TimersUsecase.GetAllByUserId(userId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted := Helper.FormatTimersForResponse(timers)
	respond.WithJson(w, http.StatusOK, timersFormatted)
}

func (u *UsersController) CreateTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	timer.Id = bson.NewObjectId()

	GetDataFromBodyRequest(r, &timer)

	err := u.TimersUsecase.Create(timer)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController) UpdateTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	GetDataFromBodyRequest(r, &timer)

	err := u.TimersUsecase.Update(timer)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController) DeleteTimer(w http.ResponseWriter, r *http.Request){

	timerId := GetIdFromUrl(r)
	err := u.TimersUsecase.Delete(timerId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController) StartTimer(w http.ResponseWriter, r *http.Request){

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
	roleUsed := Helper.WhichRoleIsUsed(userRequesting, model.User{Id: user.Id})
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

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController) FinishTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	GetDataFromBodyRequest(r, &timer)
	GetDataFromHeaderRequest(r)
	
	timerDb, err := u.TimersUsecase.GetById(timer.Id.Hex())

	// TIMER NOT EXIST
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if Helper.IsAlreadyFinished(timerDb.Finish) {
		respond.WithError(w, http.StatusBadRequest, "TimerAlreadyFinished")
		return
	}

	userIdRequesting, _ := authenticator.GetUserIdFromJWT(jwtSent)
	userRequesting, _ := u.UsersUsecase.GetById(userIdRequesting)
	fmt.Println("(FinishTimer): Id user requesting: " + userIdRequesting)
	fmt.Println("(FinishTimer): Id user wants stop: " + timerDb.UserId)

	roleUsed := Helper.WhichRoleIsUsed(userRequesting, model.User{Id: bson.ObjectIdHex(timerDb.UserId)})
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

	if err = u.TimersUsecase.Update(timerDb); err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userOwnerOfTimer, _ := u.UsersUsecase.GetById(timerDb.UserId)

	if userOwnerOfTimer.HasGroup(){
		for _, socketMap := range socketsMaps {
			for groupId, websocket := range socketMap{
				if userOwnerOfTimer.IsFromTheSameGroup(groupId) {
					message := userOwnerOfTimer.Username + " finished a timer now."
					websocket.SendMessage("response", message)
				}
			}
		}
	}

	timerFormatted := Helper.FormatTimerForResponse(timerDb)
	err = mailSender.SendFinishedTime(userOwnerOfTimer.Email, timerFormatted)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, timerFormatted)
}

// Methods to improving readability only...
func GetIdFromUrl(request *http.Request) (id string) {
	vars := mux.Vars(request)
	id = vars["id"]
	return
} // Must be abstracted?

func GetDataFromBodyRequest(r *http.Request, dataSaver interface{}){
	json.NewDecoder(r.Body).Decode(dataSaver)
}
func GetDataFromHeaderRequest(r *http.Request){
	jwtSent = Helper.GetJWTFromHeader(r)
}

// From 623 lines to 507 (04/04/19)
// From 507 lines to 537 (05/04/19)
// From 537 lines to 579 (08/04/19)
// From 579 lines to 568 (09/04/19)