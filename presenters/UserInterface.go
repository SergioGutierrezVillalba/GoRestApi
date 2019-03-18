package presenters

import (
	"FirstProject/core/usecases"
	"FirstProject/model"
	"FirstProject/model/validation"
	auth "FirstProject/model/auth"

	"fmt"
	"net/http"
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"

)

type UserInterface struct {}

var (
	userUsecase 	usecases.UserUsecase

	respond 		model.Responser
	mailSender  	model.MailSender

	crypter     	auth.Crypter
	requestHeaders 	auth.RequestHeaders

	checker 		validation.Checker
)

func (userInterface *UserInterface) GetUsers(response http.ResponseWriter, request *http.Request){

	mustDoAction, returnMsg := GetAuth(request, "GET")

	if !mustDoAction {
		respond.WithError(response, http.StatusBadRequest, returnMsg)
		return
	}

	users, err:= userUsecase.GetUsers()

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	} 

	respond.WithJson(response, http.StatusOK, users)
}

func (userInterface *UserInterface) GetUserById(response http.ResponseWriter, request *http.Request){

	mustDoAction, returnMsg := GetAuth(request, "GET")

	if !mustDoAction {
		respond.WithError(response, http.StatusBadRequest, returnMsg)
		return
	}

	userId := GetIdFromUrl(request)
	user, err:= userUsecase.GetUserById(userId);

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, "UserNotExistError")
		return
	}

	respond.WithJson(response, http.StatusOK, user)
}

func (userInterface *UserInterface) CreateUser(response http.ResponseWriter, request *http.Request){

	var user model.User
	var fieldsRequired []string

	user.Id = bson.NewObjectId()

	json.NewDecoder(request.Body).Decode(&user)	
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email)

	if !checker.HasFieldsRequired(fieldsRequired){
		respond.WithError(response, http.StatusBadRequest, "EmptyDataError")
		return
	}

	if checker.UsernameAlreadyExists(user.Username){
		respond.WithError(response, http.StatusBadRequest, "UsernameAlreadyExistsError")
		return
	}

	user.Password, _ = crypter.Crypt(user.Password)
	user.Role = "user"

	err := userUsecase.CreateUser(user);

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	} 
		
	respond.WithJson(response, http.StatusOK, "Success")
}

func (userInterface *UserInterface) UpdateUser(response http.ResponseWriter, request *http.Request){
	
	SetHeaders(request)

	var userToUpdate model.User
	var fieldsRequired []string
	jwtSent := requestHeaders.Authorization

	json.NewDecoder(request.Body).Decode(&userToUpdate)

	if userToUpdate.Id == "" {
		respond.WithError(response, http.StatusBadRequest, "EmptyIdError")
		return
	}

	fieldsRequired = append(fieldsRequired, userToUpdate.Username, userToUpdate.Password, userToUpdate.Email, jwtSent)

	if !checker.HasFieldsRequired(fieldsRequired){
		respond.WithError(response, http.StatusBadRequest, "EmptyDataError")
		return
	}

	if !checker.JwtIsCorrect(jwtSent){
		respond.WithError(response, http.StatusBadRequest, "FalseJwt")
		return
	}

	userRequesting, _ := userUsecase.GetUserByJwt(jwtSent)

	if userRequesting.Role == "user" {
		if !checker.IsUpdatingItself(userRequesting, userToUpdate){
			respond.WithError(response, http.StatusBadRequest, "UserPermissionsError")
			return
		}
	}

	userToUpdate.Password, _ = crypter.Crypt(userToUpdate.Password)
	userToUpdate.Role = "user"

	err2 := userUsecase.UpdateUser(userToUpdate)

	if err2 != nil {
		respond.WithError(response, http.StatusBadRequest, err2.Error())
		return
	} 
		
	respond.WithJson(response, http.StatusOK, "Success")
}

func (userInterface *UserInterface) DeleteUser(response http.ResponseWriter, request *http.Request){

	mustDoAction, returnMsg := GetAuth(request, "DELETE")

	if !mustDoAction {
		respond.WithError(response, http.StatusBadRequest, returnMsg)
		return
	}

	userId := GetIdFromUrl(request)
	err := userUsecase.DeleteUser(userId)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	} 
		
	respond.WithJson(response, http.StatusOK, "Success")
}

func (userInterface *UserInterface) Login(response http.ResponseWriter, request *http.Request){

	var user model.User
	json.NewDecoder(request.Body).Decode(&user)

	userDb, err := userUsecase.GetUserByUsername(user.Username)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, "UserNotExistError")
		return
	}

	err2:= crypter.PasswordCoincides(userDb.Password, user.Password)

	if err2 != nil {
		respond.WithError(response, http.StatusBadRequest, "WrongPasswordError")
		return
	}

	token := auth.GenerateJWT(user)
	result := auth.ResponseToken{token}
	userDb.Jwt = result.Token

	if err := userUsecase.UpdateUser(userDb); err != nil {
		respond.WithError(response, http.StatusBadRequest, "UpdateJwtError")
	}

	respond.WithJson(response, http.StatusOK, result)
}

func (userInterface *UserInterface) SendRecover(response http.ResponseWriter, request *http.Request){

	var user model.User
	json.NewDecoder(request.Body).Decode(&user)

	if !checker.HasFieldsRequired([]string{user.Username}) {
		respond.WithError(response, http.StatusBadRequest, "EmptyDataError")
		return
	}

	if err:= checker.CheckSpecialChars(user.Username); err != nil {
		respond.WithError(response, http.StatusBadRequest, "Invalid chars")
		return
	} 

	if !checker.UsernameAlreadyExists(user.Username) {
		respond.WithError(response, http.StatusBadRequest, "UserNotExistError")
		return
	} 

	if email, err:= GetEmailUser(user.Username); err != nil {
		respond.WithError(response, http.StatusBadRequest, "EmailError")
		return

	} else {

		token, err := GenerateToken(email)

		if err != nil {
			respond.WithError(response, http.StatusBadRequest, "TokenGeneratorError")
			return
		}

		err = mailSender.Send(email, token)

		if err != nil {
			respond.WithError(response, http.StatusBadRequest, "SendingError")
			return
		}	

		err = SaveToken(token, user.Username)

		if err != nil {
			respond.WithError(response, http.StatusBadRequest, "SavingTokenError")
			return
		}

	}

	respond.WithJson(response, http.StatusOK, "Success")
}

func (userInterface *UserInterface) Reset(response http.ResponseWriter, request *http.Request){

	var passwordRecover model.PasswordRecover
	var fieldsRequired []string

	json.NewDecoder(request.Body).Decode(&passwordRecover)

	fieldsRequired = append(fieldsRequired, passwordRecover.Token, passwordRecover.NewPassword)

	if !checker.HasFieldsRequired(fieldsRequired){
		respond.WithError(response, http.StatusBadRequest, "EmptyDataError")
		return
	}

	user, err := userUsecase.GetUserByToken(passwordRecover.Token);

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, "FalseTokenError")
		return
	}

	user.Password, err = crypter.Crypt(passwordRecover.NewPassword)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, "CrypterError")
		return
	}

	err = userUsecase.UpdateUser(user)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, "UpdateUserError")
		return
	}

	respond.WithJson(response, http.StatusOK, "Success")
}


// Func
func SetHeaders(request *http.Request){
	requestHeaders.SetHeaders(request)
}

func GetIdFromUrl(request *http.Request) string {

	vars := mux.Vars(request)
	userId := vars["id"]

	return userId
}

func GetEmailUser(username string)(string, error){

	user, err:= userUsecase.GetEmailOfUser(username)
	return user.Email, err
}

func GenerateToken(hash string) (string, error){

	token, err := crypter.Crypt(hash)
	return token, err
}

func SaveToken(token string, username string) error {
	err:= userUsecase.SaveToken(token, username)
	return err
}

func GetAuth(request *http.Request, httpMethodRequested string) (doAction bool, msg string) {

	SetHeaders(request)

	doAction = true
	msg = "Success"
	jwtSent := requestHeaders.Authorization


	if !checker.HasFieldsRequired([]string {jwtSent}){
		fmt.Println("Aquí estoy, linea 320")
		msg = "EmptyJwtError"
		doAction = false
		return
	}

	fmt.Println("JWT: " + jwtSent)
	fmt.Println("Aqui estoy, linea 339")
	if !checker.JwtIsCorrect(jwtSent){
		fmt.Println("Aquí estoy, linea 327")
		msg = "JwtError"
		doAction = false
		return
	}

	userRequesting, _ := userUsecase.GetUserByJwt(jwtSent)

	if !checker.HasPermissions(userRequesting.Role, httpMethodRequested){
		fmt.Println("Aquí estoy, linea 336")
		msg = "UserPermissionsError"
		doAction = false
		return
	}

	return
}