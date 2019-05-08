package controller

import (
	timersUsecase "FirstProject/Domains/timer/usecase"
	usersUsecase "FirstProject/Domains/user/usecase"
	"errors"

	model "FirstProject/Model"
	auth "FirstProject/Model/Auth"
	"FirstProject/Model/helper"
	"FirstProject/Model/imgs"
	"FirstProject/Model/mail"
	"FirstProject/Model/socket"

	swagger "FirstProject/Swagger"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/mgo.v2/bson"
)

type Controller interface {
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
	UsersUsecase  usersUsecase.Usecase
	TimersUsecase timersUsecase.Usecase
}

var (
	socketsMaps = make(map[string]map[string]*socket.WebSocket)

	respond model.Responser

	mailSender mail.MailSender

	Helper helper.Helper

	jwtSent string

	crypter       auth.Crypter
	authenticator auth.Authentication
	responseToken auth.ResponseToken
	requestInfo   auth.RequestInfo
)

// NewController generates a new controller
// Params: Usecases
func NewController(u usersUsecase.Usecase, t timersUsecase.Usecase) Controller {
	return &UsersController{
		UsersUsecase:  u,
		TimersUsecase: t,
	}
}

// SOCKETS CONTEXT

// StartWebSocket
func (u *UsersController) StartWebSocket(w http.ResponseWriter, r *http.Request) {

	ws, err := socket.NewWebSocket(w, r)
	var tokenReceived string

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO SocketsList struct
	ws.On("message", func(e *socket.Event) {
		tokenReceived = e.Data.(string)
		uncryptedJWT := authenticator.Decrypt(authenticator.DecodeBase64(tokenReceived))
		userFromJWT := authenticator.GetUserInfoFromJWT(uncryptedJWT)
		fmt.Println("(StartWebSocket): GroupId - " + userFromJWT.GroupId)

		usersGroupId := userFromJWT.GroupId
		socketMap := CreateSocketMap(ws, usersGroupId)
		AddSocketMapToSocketsMaps(socketMap)
		// idMap := ksuid.New().String()
		// groupAndSocket := make(map[string]*socket.WebSocket)
		// groupAndSocket[userFromJWT.GroupId] = ws
		// socketsMaps[idMap] = groupAndSocket

		// log.Printf("[MESSAGE] %v", e.Data)
		ws.Out <- (&socket.Event{
			Name: "response",
			Data: "Socket created",
		}).Raw()
	})
}

// CreateSocketMap
func CreateSocketMap(ws *socket.WebSocket, groupId string) map[string]*socket.WebSocket {
	socketMap := make(map[string]*socket.WebSocket)
	socketMap[groupId] = ws
	return socketMap
}

// AddSocketMapToSocketsMap
func AddSocketMapToSocketsMaps(socketMap map[string]*socket.WebSocket) {
	randomIdMap := ksuid.New().String()
	socketsMaps[randomIdMap] = socketMap
}

// FinishWebSocket
func (u *UsersController) FinishWebSocket(w http.ResponseWriter, r *http.Request) {

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

// GetAllUsers swagger:route GET /users users getAllUsers
// Returns all users from Database.
// Requires Auth (JWT)
// responses:
//   200: body:getAllUsersResp
//   400: body:queryErrResp
func (u *UsersController) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := u.UsersUsecase.GetAll()

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, swagger.GenericQueryErrResp{
			Error: err,
		})
		return
	}

	respond.WithJson(w, http.StatusOK, swagger.GetAllUsersResp{Users: users})
}

// GetUserByJwt swagger:route GET /user users getUserByJwt
// Returns user by JWT sent.
// Requires Auth (JWT)
// responses:
//  200: body:getUserResp
//  400: body:queryErrResp
func (u *UsersController) GetUserByJwt(w http.ResponseWriter, r *http.Request) {

	GetDataFromHeaderRequest(r)
	user, err := u.UsersUsecase.GetUserByJwt(jwtSent)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, swagger.GenericQueryErrResp{
			Error: err,
		})
		return
	}

	respond.WithJson(w, http.StatusOK, swagger.GetUserResp{User: user})
}

// GetUserById swagger:operation GET /users/{id} users getUserById
// ---
// summary: Returns an user by Id sent
// description: if user id is not send correctly 400 code is returned
// parameters:
// - name: id
//   in: path
//   description: id of a user
//   type: string
//   required: true
// responses:
//  200: body:getUserResp
//  400: body:queryErrResp
func (u *UsersController) GetUserById(w http.ResponseWriter, r *http.Request) {

	userId := GetIdFromUrl(r)
	user, err := u.UsersUsecase.GetById(userId)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, swagger.GenericQueryErrResp{
			Error: err,
		})
		return
	}

	respond.WithJson(w, http.StatusOK, swagger.GetUserResp{User: user})
}

// CreateUser swagger:route POST /users users createUser
// Creates a user.
// responses:
//  200: body:genericSuccessResp
//  400: body:queryErrResp
func (u *UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user model.User
	user.Id = bson.NewObjectId()

	GetDataFromBodyJSONRequest(r, &user)

	user.Password, _ = crypter.Crypt(user.Password)

	if ActionGivesError(u.UsersUsecase.Create(user)) {
		respond.WithError(w, http.StatusBadRequest, swagger.GenericQueryErrResp{
			Error: errors.New("CreationUserError"),
		})
		return
	}

	user.EmptyPassword()
	respond.WithJson(w, http.StatusOK, swagger.GenericSuccessResp{Response: "Success"})
}

// UpdateUser swagger:route PUT /users users updateUser
// Updates a user.
// Updates a user with JSON given and returns updated JWT.
// responses:
//  200: body:updateUserResp
//  400: body:queryErrResp
func (u *UsersController) UpdateUser(w http.ResponseWriter, r *http.Request) {

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

	if ActionGivesError(u.UsersUsecase.Update(userToUpdate)) {
		respond.WithError(w, http.StatusBadRequest, "UpdateUserError")
		return
	}
	respond.WithJson(w, http.StatusOK, auth.ResponseToken{Jwt: userToUpdate.Jwt})
}

// UpdateUserWithoutUpdatingPassword swagger:route PUT /users/nopwd users updateUserWithoutPwd
// Updates a user without updating password.
// Updates a user with JSON given without changing password
// responses:
//  200: body:updateUserResp
//  400: body:queryErrResp
func (u *UsersController) UpdateUserWithoutUpdatingPassword(w http.ResponseWriter, r *http.Request) {

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
	respond.WithJson(w, http.StatusOK, auth.ResponseToken{Jwt: userToUpdate.Jwt})
}

// DeleteUser swagger:route DELETE /users/{id} users deleteUser
// Delete a user with given Id.
// Deletes a user with given Id via url
// responses:
//  200: body:genericSuccessResp
//  400: body:queryErrResp
func (u *UsersController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	idUser := GetIdFromUrl(r)
	err := u.UsersUsecase.Delete(idUser)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, struct {
		Response string `json:"response"`
	}{
		"Success",
	})
}

// Login swagger:route POST /users/login validation login
// Login.
// Checks data sent for login and returns jwt if successful
// responses:
//  200: body:loginResp
//  400: body:queryErrResp
func (u *UsersController) Login(w http.ResponseWriter, r *http.Request) {

	var user model.User

	GetDataFromBodyJSONRequest(r, &user)
	userDb, err := u.UsersUsecase.GetUserByUsername(user.Username)

	if err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if ActionGivesError(u.UsersUsecase.Login(user, userDb)) {
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
	respond.WithJson(w, http.StatusOK, auth.ResponseToken{Jwt: userDb.Jwt})
}

// Register swagger:route POST /users/register validation register
// Register.
// Allows register for new users
// responses:
//  200: body:genericSuccessResp
//  400: body:queryErrResp
func (u *UsersController) Register(w http.ResponseWriter, r *http.Request) {

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

// SendRecover swagger:route POST /users/sendrecover validation sendRecover
// Allows way to recover password when user has forgotten it.
// Sends an mail to user's email related to with an url to recover password
// responses:
//  200: body:genericSuccessResp
//  400: body:queryErrResp
func (u *UsersController) SendRecover(w http.ResponseWriter, r *http.Request) {

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

	token = CleanSlashesFromToken(token)
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

// ResetPassword swagger:route PATCH /users/reset validation resetPassword
// Change password.
// Change password with recover token given.
// responses:
//  200: body:genericSuccessResp
//  400: body:queryErrResp
func (u *UsersController) ResetPassword(w http.ResponseWriter, r *http.Request) {

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

	if err = u.UsersUsecase.Update(user); err != nil {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}

// SetProfileImage
func (u *UsersController) SetProfileImage(w http.ResponseWriter, r *http.Request) {

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

// GetProfileImage
func (u *UsersController) GetProfileImage(w http.ResponseWriter, r *http.Request) {

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

// TIMERS CONTEXT

// GetAllTimers
func (u *UsersController) GetAllTimers(w http.ResponseWriter, r *http.Request) {

	timers, err := u.TimersUsecase.GetAll()

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted := Helper.FormatTimersForResponse(timers)
	respond.WithJson(w, http.StatusOK, timersFormatted)
}

// GetTimerById
func (u *UsersController) GetTimerById(w http.ResponseWriter, r *http.Request) {

	timerId := GetIdFromUrl(r)
	timer, err := u.TimersUsecase.GetById(timerId)

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

// GetTimersByUserId
func (u *UsersController) GetTimersByUserId(w http.ResponseWriter, r *http.Request) {

	userId := GetIdFromUrl(r)
	timers, err := u.TimersUsecase.GetAllByUserId(userId)

	if ActionGivesError(err) {
		respond.WithError(w, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted := Helper.FormatTimersForResponse(timers)
	respond.WithJson(w, http.StatusOK, timersFormatted)
}

// CreateTimer
func (u *UsersController) CreateTimer(w http.ResponseWriter, r *http.Request) {

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

// UpdateTimer
func (u *UsersController) UpdateTimer(w http.ResponseWriter, r *http.Request) {

	var timer model.Timer
	GetDataFromBodyJSONRequest(r, &timer)

	if ActionGivesError(u.TimersUsecase.Update(timer)) {
		respond.WithError(w, http.StatusBadRequest, "UpdateTimerError")
		return
	}

	timerFormatted := Helper.FormatTimerForResponse(timer)
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

// DeleteTimer
func (u *UsersController) DeleteTimer(w http.ResponseWriter, r *http.Request) {

	timerId := GetIdFromUrl(r)

	if ActionGivesError(u.TimersUsecase.Delete(timerId)) {
		respond.WithError(w, http.StatusBadRequest, "DeleteTimerError")
		return
	}

	respond.WithJson(w, http.StatusOK, "Success")
}

// StartTimer
func (u *UsersController) StartTimer(w http.ResponseWriter, r *http.Request) {

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

// FinishTimer
func (u *UsersController) FinishTimer(w http.ResponseWriter, r *http.Request) {

	var timer model.Timer
	GetDataFromBodyJSONRequest(r, &timer)
	GetDataFromHeaderRequest(r)

	timerDb, _ := u.TimersUsecase.GetById(timer.GetId())

	if timerDb.IsAlreadyFinished() {
		respond.WithError(w, http.StatusBadRequest, "TimerAlreadyFinished")
	}

	userRequesting := u.GetUserRequesting()
	fmt.Println("(FinishTimer): Id user wants stop: " + timerDb.UserId)

	roleUsed := Helper.WhichRoleIsUsed(userRequesting, model.User{Id: timerDb.UserId})
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

	userOwnerOfTimer, _ := u.UsersUsecase.GetById(timerDb.UserId.Hex())

	if userOwnerOfTimer.HasGroup() {
		SendFinishNotificationToTheGroup(userOwnerOfTimer)
	}

	timerFormatted := FormatTimerForResponse(timerDb)
	if ActionGivesError(mailSender.SendFinishedTime(userOwnerOfTimer.Email, timerFormatted)) {
		respond.WithError(w, http.StatusBadRequest, "SendingEmailError")
		return
	}
	respond.WithJson(w, http.StatusOK, timerFormatted)
}

// GetIdFromUrl
func GetIdFromUrl(r *http.Request) (id string) {
	id = Helper.GetIdFromUrl(r)
	return
}

// GetDataFromBodyJSONRequest
func GetDataFromBodyJSONRequest(r *http.Request, dataSaver interface{}) {
	json.NewDecoder(r.Body).Decode(dataSaver)
}

// GetDataFromHeaderRequest
func GetDataFromHeaderRequest(r *http.Request) {
	jwtSent, _ = Helper.GetJWTFromHeaderRequest(r)
}

// SendFinishNotificationToTheGroup
func SendFinishNotificationToTheGroup(userOwnerOfTimer model.User) {
	for _, socketMap := range socketsMaps {
		for groupId, websocket := range socketMap {
			if userOwnerOfTimer.IsFromTheSameGroup(groupId) {
				message := userOwnerOfTimer.Username + " finished a timer now."
				websocket.SendMessage("response", message)
			}
		}
	}
}

// GetUserRequesting
func (u *UsersController) GetUserRequesting() model.User {
	userIdRequesting, _ := authenticator.GetUserIdFromJWT(jwtSent)
	userRequesting, _ := u.UsersUsecase.GetById(userIdRequesting)
	fmt.Println("(GetUserRequesting): Id user requesting: " + userIdRequesting)
	return userRequesting
}

// GenerateJWTAndSaveInUser
// Needs a pointer for saving JWT everywhere it's called.
// At the same time, needs the struct for generating a JWT.
func GenerateJWTAndSaveInUser(userPointer *model.User) {
	var user model.User
	user = *userPointer
	newJWT := authenticator.GenerateJWT(user)
	userPointer.SetJWT(newJWT)
}

// SaveFinishAndDuration
func SaveFinishAndDuration(timerPointer *model.Timer) {
	finishTime := time.Now().Unix()
	duration := finishTime - timerPointer.Start

	timerPointer.Finish = finishTime
	timerPointer.Duration = duration
}

// FormatTimerForResponse
func FormatTimerForResponse(timer model.Timer) model.TimerFormatted {
	return Helper.FormatTimerForResponse(timer)
}

// ActionGivesError
func ActionGivesError(e error) bool {
	return Helper.ActionGivesError(e)
}

// CleanUserPasswordAndJWT
func CleanUserPasswordAndJWT(userPointer *model.User) {
	userPointer.EmptyPassword()
	userPointer.EmptyJWT()
}

// CreateTimerStruct
func CreateTimerStruct(userOwner model.User) (timer model.Timer) {
	timer.Id = bson.NewObjectId()
	timer.UserId = userOwner.Id
	timer.Start = time.Now().Unix()
	return
}

// PrepareUserForUpdate
func PrepareUserForUpdate(userPointer *model.User, newPassword string) {
	userPointer.EmptyPassword()
	GenerateJWTAndSaveInUser(userPointer)
	userPointer.SetPassword(newPassword)
}

// CleanSlashesFromToken
func CleanSlashesFromToken(token string) string {
	return Helper.CleanSlashesFromToken(token)
}
