package controllers

import (
	"FirstProject/core/usecases"
	"FirstProject/model"

	// "fmt"
	"net/http"
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"

)

type UserInterface struct {}

var (
	userUsecase 	usecases.UserUsecase
	jwtSent			string
)

func (userInterface *UserInterface) GetUsers(response http.ResponseWriter, request *http.Request){

	GetDataFromRequest(request, nil)

	users, err:= userUsecase.GetUsers(jwtSent)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	} 

	respond.WithJson(response, http.StatusOK, users)
}

func (userInterface *UserInterface) GetUserById(response http.ResponseWriter, request *http.Request){

	GetDataFromRequest(request, nil)

	userId := GetIdFromUrl(request)
	user, err:= userUsecase.GetUserById(jwtSent, userId);

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(response, http.StatusOK, user)
}

func (userInterface *UserInterface) CreateUser(response http.ResponseWriter, request *http.Request){

	var user model.User

	user.Id = bson.NewObjectId()

	GetDataFromRequest(request, &user)

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

	var userToUpdate model.User

	GetDataFromRequest(request, &userToUpdate)

	userToUpdate.Password, _ = crypter.Crypt(userToUpdate.Password)
	userToUpdate.Role = "user"

	err := userUsecase.UpdateUser(jwtSent, userToUpdate)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	} 
		
	respond.WithJson(response, http.StatusOK, "Success")
}

func (userInterface *UserInterface) DeleteUser(response http.ResponseWriter, request *http.Request){

	GetDataFromRequest(request, nil)

	userId := GetIdFromUrl(request)
	err := userUsecase.DeleteUser(jwtSent, userId)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	} 
		
	respond.WithJson(response, http.StatusOK, "Success")
}

func (userInterface *UserInterface) Login(response http.ResponseWriter, request *http.Request){

	var user model.User
	
	GetDataFromRequest(request, &user)

	userDb, err := userUsecase.GetUserByUsername(user.Username)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	err2:= crypter.PasswordCoincides(userDb.Password, user.Password)

	if err2 != nil {
		respond.WithError(response, http.StatusBadRequest, err2.Error())
		return
	}

	token := authenticator.GenerateJWT(user)
	responseToken.Token = token
	userDb.Jwt = responseToken.Token

	if err := userUsecase.UpdateUserWithoutJwt(userDb); err != nil {
		respond.WithError(response, http.StatusBadRequest, "UpdateJwtError")
		return
	}

	respond.WithJson(response, http.StatusOK, responseToken)
}

func (userInterface *UserInterface) SendRecover(response http.ResponseWriter, request *http.Request){

	var user model.User
	
	GetDataFromRequest(request, &user)

	if email, err:= GetEmailUser(user.Username); err != nil {
		respond.WithError(response, http.StatusBadRequest, "UserNotExistError")
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

	json.NewDecoder(request.Body).Decode(&passwordRecover)

	user, err := userUsecase.GetUserByToken(passwordRecover.Token);

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	user.Password, err = crypter.Crypt(passwordRecover.NewPassword)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	err = userUsecase.UpdateUserWithoutJwt(user)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(response, http.StatusOK, "Success")
}


// Func

func GetDataFromRequest(request *http.Request, dataSaver interface{}){
	SetHeaders(request)
	GetJwtSent()
	json.NewDecoder(request.Body).Decode(dataSaver)
}

func GetJwtSent() {
	jwtSent = requestHeaders.Authorization
}

func GetEmailUser(username string)(string, error){
	user, err:= userUsecase.GetEmailOfUser(username)
	return user.Email, err
}

func GetIdFromUrl(request *http.Request) string {

	vars := mux.Vars(request)
	userId := vars["id"]

	return userId
}

func SetHeaders(request *http.Request){
	requestHeaders.SetHeaders(request)
}


// Tokenizer abstraction ?? 
func GenerateToken(hash string) (string, error){
	token, err := crypter.Crypt(hash)
	return token, err
}

func SaveToken(token string, username string) error {
	err:= userUsecase.SaveToken(token, username)
	return err
}