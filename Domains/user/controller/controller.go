package controller

import (
	usersUsecase "FirstProject/Domains/user/usecase"
	timersUsecase "FirstProject/Domains/timer/usecase"

	"FirstProject/Model/Auth"
	"FirstProject/Model/helper"
	"FirstProject/Model/socket"
	"FirstProject/Model/mail"
	"FirstProject/Model/imgs"
	"FirstProject/Model"

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
	UpdateUserWithoutUpdatingPassword(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	SendRecover(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	SetProfileImage(w http.ResponseWriter, r *http.Request)
	GetProfileImage(w http.ResponseWriter, r *http.Request)
	
	StartWebSocket(w http.ResponseWriter, r *http.Request)
	FinishWebSocket(w http.ResponseWriter, r *http.Request)

	GetTasksOnTheSameDateAsUserTimersByUserId(w http.ResponseWriter, r *http.Request)

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
	GetDataFromBodyJSONRequest(r, &socket)
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

	user, err := u.UsersUsecase.GetUserByJwt(jwtSent)
	user.EmptyPassword()
	user.EmptyJWT()

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

	GetDataFromBodyJSONRequest(r, &user)

	user.Password, _ = crypter.Crypt(user.Password)

	if ActionGivesError(u.UsersUsecase.Create(user)) {
		respond.WithError(w, http.StatusBadRequest, "CreateUserError")
		return
	}

	user.EmptyPassword()
	respond.WithJson(w, http.StatusOK, user)
}

func (u *UsersController) UpdateUser(w http.ResponseWriter, r *http.Request){

	var userToUpdate model.User
	GetDataFromBodyJSONRequest(r, &userToUpdate)
	GetDataFromHeaderRequest(r)

	userRequesting := u.GetUserRequesting()
	fmt.Println("(UpdateUser): Id user wants upda: " + userToUpdate.Id.Hex())

	roleUsed := Helper.WhichRoleIsUsed(userRequesting, userToUpdate)
	switch roleUsed {
		case "NOAUTH":
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		case "SELF":
			userToUpdate.SetRole("user")
		case "ADMIN":
			// add mongodb parameter to not allow changing password
	}

	newPassword, _ := crypter.Crypt(userToUpdate.Password)
	PrepareUserForUpdate(&userToUpdate, newPassword)

	if ActionGivesError(u.UsersUsecase.Update(userToUpdate)){
		respond.WithError(w, http.StatusBadRequest, "UpdateUserError")
		return
	}
	respond.WithJson(w, http.StatusOK, auth.ResponseToken{Token:userToUpdate.Jwt})
}

// TODO switch can be replaced to not doing it 
// twice each time I try to check the userRole?
func (u *UsersController) UpdateUserWithoutUpdatingPassword(w http.ResponseWriter, r *http.Request){

	var userToUpdate model.User
	GetDataFromBodyJSONRequest(r, &userToUpdate)
	GetDataFromHeaderRequest(r)

	userRequesting := u.GetUserRequesting()
	fmt.Println("(UpdateUser): Id user wants upda: " + userToUpdate.GetId())

	userToUpdateInBD, _ := u.UsersUsecase.GetById(userToUpdate.GetId())

	roleUsed := Helper.WhichRoleIsUsed(userRequesting, userToUpdateInBD)
	switch roleUsed {
		case "NOAUTH":
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		case "SELF":
			userToUpdate.SetRole("user")
		case "ADMIN":
	}

	newPassword := userToUpdateInBD.Password
	PrepareUserForUpdate(&userToUpdate, newPassword)

	if ActionGivesError(u.UsersUsecase.Update(userToUpdate)) {
		respond.WithError(w, http.StatusBadRequest, "UpdateUserError")
		return
	}
	respond.WithJson(w, http.StatusOK, auth.ResponseToken{ Token:userToUpdate.Jwt })	
}

func (u *UsersController) DeleteUser(w http.ResponseWriter, r *http.Request){

	idUser := GetIdFromUrl(r)
	err := u.UsersUsecase.Delete(idUser);

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController) Login(w http.ResponseWriter, r *http.Request){

	var user model.User

	GetDataFromBodyJSONRequest(r, &user)
	userDb, err := u.UsersUsecase.GetUserByUsername(user.Username)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if ActionGivesError(u.UsersUsecase.Login(user, userDb)){
		respond.WithError(w, http.StatusBadRequest, "LoginError")
		return
	}

	savedPWD := userDb.Password
	CleanUserPasswordAndJWT(&userDb)
	GenerateJWTAndSaveInUser(&userDb)
	userDb.SetPassword(savedPWD)

	if ActionGivesError(u.UsersUsecase.Update(userDb)) {
		respond.WithError(w, http.StatusBadRequest, "UpdateUserError")
		return
	}
	respond.WithJson(w, http.StatusOK, auth.ResponseToken{Token:userDb.Jwt})
}

func (u *UsersController) Register(w http.ResponseWriter, r *http.Request){

	var user model.User
	user.Id = bson.NewObjectId()

	GetDataFromBodyJSONRequest(r, &user)

	newPassword, _ := crypter.Crypt(user.Password)

	user.SetPassword(newPassword)
	user.SetRole("user")

	if ActionGivesError(u.UsersUsecase.Register(user)) {
		respond.WithError(w, http.StatusBadRequest, "RegisterError")
		return
	}
	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController) SendRecover(w http.ResponseWriter, r *http.Request){
	
	// REFACTOR
	var user model.User
	GetDataFromBodyJSONRequest(r, &user)

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
	GetDataFromBodyJSONRequest(r, &passwordRecover)

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

func (u *UsersController) SetProfileImage(w http.ResponseWriter, r *http.Request){

	userId := r.FormValue("id")
	multiPartFile, _, err := r.FormFile("img")

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if ActionGivesError(u.UsersUsecase.SetProfileImage(userId, multiPartFile)) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController) GetProfileImage(w http.ResponseWriter, r *http.Request){
	
	userId := GetIdFromUrl(r)
	imageBytes, err := u.UsersUsecase.GetProfileImage(userId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, imgs.ProfileImage{
		ImageBytes: imageBytes,
	})
}

// TASKS CONTEXT

func (u *UsersController) GetTasksOnTheSameDateAsUserTimersByUserId(w http.ResponseWriter, r *http.Request){

	userId := GetIdFromUrl(r)
	tasks, err := u.UsersUsecase.GetTasksOnTheSameDateAsUserTimers(userId)

	if ActionGivesError(err){
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respond.WithJson(w, http.StatusOK, tasks)
}

// TIMERS CONTEXT

func (u *UsersController) GetAllTimers(w http.ResponseWriter, r *http.Request){
	
	timers, err := u.TimersUsecase.GetAll()

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted := Helper.FormatTimersForResponse(timers)
	respond.WithJson(w, http.StatusOK, timersFormatted)
}

func (u *UsersController) GetTimerById(w http.ResponseWriter, r *http.Request){

	timerId := GetIdFromUrl(r)
	timer, err := u.TimersUsecase.GetById(timerId)

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController) GetTimersByUserId(w http.ResponseWriter, r *http.Request){

	userId := GetIdFromUrl(r)
	timers, err := u.TimersUsecase.GetAllByUserId(userId)

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted := Helper.FormatTimersForResponse(timers)
	respond.WithJson(w, http.StatusOK, timersFormatted)
}

func (u *UsersController) CreateTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	timer.Id = bson.NewObjectId()

	GetDataFromBodyJSONRequest(r, &timer)

	if ActionGivesError(u.TimersUsecase.Create(timer)) {
		respond.WithError(w, http.StatusBadRequest, "CreateTimerError")
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController) UpdateTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	GetDataFromBodyJSONRequest(r, &timer)

	if ActionGivesError(u.TimersUsecase.Update(timer)) {
		respond.WithError(w, http.StatusBadRequest, "UpdateTimerError")
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController) DeleteTimer(w http.ResponseWriter, r *http.Request){

	timerId := GetIdFromUrl(r)
	
	if ActionGivesError(u.TimersUsecase.Delete(timerId)) {
		respond.WithError(w, http.StatusBadRequest, "DeleteTimerError")
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}

func (u *UsersController) StartTimer(w http.ResponseWriter, r *http.Request){

	var user model.User

	GetDataFromBodyJSONRequest(r, &user)
	GetDataFromHeaderRequest(r)

	userRequesting := u.GetUserRequesting()
	fmt.Println("(StartTimer): Id user wants to start timer: " + userRequesting.GetId())

	roleUsed := Helper.WhichRoleIsUsed(userRequesting, model.User{Id: user.Id})
	switch roleUsed {
		case "NOAUTH":
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		case "SELF":
		case "ADMIN":
	}

	timer := CreateTimerStruct(user)

	if ActionGivesError(u.TimersUsecase.StartTimer(timer)) {
		respond.WithError(w, http.StatusBadRequest, "StartTimerError")
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

func (u *UsersController) FinishTimer(w http.ResponseWriter, r *http.Request){

	var timer model.Timer
	GetDataFromBodyJSONRequest(r, &timer)
	GetDataFromHeaderRequest(r)
	
	timerDb, _ := u.TimersUsecase.GetById(timer.GetId())

	if timerDb.IsAlreadyFinished(){
		respond.WithError(w, http.StatusBadRequest, "TimerAlreadyFinished")
	}

	userRequesting := u.GetUserRequesting()
	fmt.Println("(FinishTimer): Id user wants stop: " + timerDb.UserId)

	roleUsed := Helper.WhichRoleIsUsed(userRequesting, model.User{Id: bson.ObjectIdHex(timerDb.UserId)})
	switch roleUsed {
		case "NOAUTH":
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		case "SELF":
		case "ADMIN":
	}	

	SaveFinishAndDuration(&timerDb)

	if ActionGivesError(u.TimersUsecase.FinishTimer(timerDb)) {
		respond.WithError(w, http.StatusBadRequest, "FinishTimerError")
		return
	}

	userOwnerOfTimer, _ := u.UsersUsecase.GetById(timerDb.UserId)

	if userOwnerOfTimer.HasGroup(){
		SendFinishNotificationToTheGroup(userOwnerOfTimer)
	}

	timerFormatted := FormatTimerForResponse(timerDb)
	if ActionGivesError(mailSender.SendFinishedTime(userOwnerOfTimer.Email, timerFormatted)) {
		respond.WithError(w, http.StatusBadRequest, "SendingEmailError")
		return
	}
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

//
func GetIdFromUrl(request *http.Request) (id string) {
	vars := mux.Vars(request)
	id = vars["id"]
	return
}
func GetDataFromBodyJSONRequest(r *http.Request, dataSaver interface{}){
	json.NewDecoder(r.Body).Decode(dataSaver)
}
func GetDataFromHeaderRequest(r *http.Request){
	jwtSent = Helper.GetJWTFromHeader(r)
}
func SendFinishNotificationToTheGroup(userOwnerOfTimer model.User){
	for _, socketMap := range socketsMaps {
		for groupId, websocket := range socketMap{
			if userOwnerOfTimer.IsFromTheSameGroup(groupId) {
				message := userOwnerOfTimer.Username + " finished a timer now."
				websocket.SendMessage("response", message)
			}
		}
	}
}
func (u *UsersController) GetUserRequesting() model.User {
	userIdRequesting, _ := authenticator.GetUserIdFromJWT(jwtSent) 
	userRequesting, _ := u.UsersUsecase.GetById(userIdRequesting)
	fmt.Println("(GetUserRequesting): Id user requesting: " + userIdRequesting)
	return userRequesting
}

// Needs a pointer for saving JWT everywhere it's called.
// At the same time, needs the struct for generating a JWT.

func GenerateJWTAndSaveInUser(userPointer *model.User){
	var user model.User
	user = *userPointer
	newJWT := authenticator.GenerateJWT(user)
	userPointer.SetJWT(newJWT)
}
func SaveFinishAndDuration(timerPointer *model.Timer){
	finishTime := time.Now().Unix()
	duration := finishTime - timerPointer.Start

	timerPointer.Finish = finishTime
	timerPointer.Duration = duration
}
func FormatTimerForResponse(timer model.Timer) model.TimerFormatted {
	return Helper.FormatTimerForResponse(timer)
}
func ActionGivesError(err error) bool {
	if err != nil {
		return true
	}
	return false
}
func CleanUserPasswordAndJWT(userPointer *model.User){
	userPointer.EmptyPassword()
	userPointer.EmptyJWT()
}
func CreateTimerStruct(userOwner model.User) (timer model.Timer) {
	timer.Id = bson.NewObjectId()
	timer.UserId = userOwner.GetId()
	timer.Start = time.Now().Unix()
	return
}
func PrepareUserForUpdate(userPointer *model.User, newPassword string){
	userPointer.EmptyPassword()
	GenerateJWTAndSaveInUser(userPointer)
	userPointer.SetPassword(newPassword)
}