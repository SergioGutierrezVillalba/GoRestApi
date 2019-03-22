package usecases

import (
	"FirstProject/model"

	"errors"
	"fmt"
	// "regexp"
)

type UserUsecase struct {
	usecase Usecase
}


func (userUsecase *UserUsecase) GetUsers(jwtSent string) ([]model.User, error){

	mustDoAction, returnMsg := userUsecase.GetAuth(jwtSent, "GET")

	if !mustDoAction {
		return nil, errors.New(returnMsg)
	}

	return repository.GetUsers()
}

func (userUsecase *UserUsecase) GetUserById(jwtSent string, userId string) (interface{}, error){

	mustDoAction, returnMsg := userUsecase.GetAuth(jwtSent, "GET")

	if !mustDoAction {
		return nil, errors.New(returnMsg)
	}

	return repository.GetUserById(userId)
}

func (userUsecase *UserUsecase) GetUserByUsername(username string) (model.User, error){

	var user model.User

	if username == "" {
		return user, errors.New("UsernameEmptyError")
	}

	return repository.GetUserByUsername(username)
}

func (userUsecase *UserUsecase) GetUserByJwt(jwt string)(model.User, error){
	return repository.GetUserByJwt(jwt)
}

func (userUsecase *UserUsecase) GetEmailOfUser(username string) (model.User, error){
	return repository.GetEmailOfUser(username)
}

func (userUsecase *UserUsecase) GetUserByToken(token string) (model.User, error){

	var user model.User

	if token == "" {
		return user, errors.New("EmptyTokenError")
	}

	return repository.GetUserByToken(token)
}


func (userUsecase *UserUsecase) FindJwt(jwtReceived string) error{
	return repository.FindJwt(jwtReceived)
}

func (userUsecase *UserUsecase) CreateUser(user model.User) error{
	
	var fieldsRequired []string

	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email)

	if !checker.HasFieldsRequired(fieldsRequired){
		return errors.New("EmptyDataError")
	}

	_, err := userUsecase.GetUserByUsername(user.Username)

	if err == nil {
		return errors.New("UsernameAlreadyExistsError")
	}

	return repository.CreateUser(user)
}

func (userUsecase *UserUsecase) UpdateUser(jwtReceived string, userToUpdate model.User) error{

	var fieldsRequired []string

	if userToUpdate.Id == "" {
		return errors.New("EmptyIdError")
	}

	fieldsRequired = append(fieldsRequired, userToUpdate.Username, userToUpdate.Password, userToUpdate.Email, jwtReceived)

	if !checker.HasFieldsRequired(fieldsRequired){
		return errors.New("EmptyDataError")
	}

	userRequesting, err := userUsecase.GetUserByJwt(jwtReceived)

	if err != nil {
		return errors.New("InvalidJwtError")
	}

	if userRequesting.Role == "user" {
		if !checker.IsUpdatingItself(userRequesting, userToUpdate){
			return errors.New("UserPermissionsError")
		}
	}

	return repository.UpdateUser(userToUpdate)
}

func (userUsecase *UserUsecase) UpdateUserWithoutJwt(userToUpdate model.User) error {

	if !checker.HasFieldsRequired([]string{userToUpdate.Password}){
		return errors.New("EmptyDataError")
	}

	return repository.UpdateUser(userToUpdate)
}

func (userUsecase *UserUsecase) DeleteUser(jwtReceived string, userId string) error{

	mustDoAction, returnMsg := userUsecase.GetAuth(jwtReceived, "DELETE")

	if !mustDoAction {
		return errors.New(returnMsg)
	}

	return repository.DeleteUser(userId)
}

func (userUsecase *UserUsecase) SaveToken(token string, username string) error {
	user, _ := userUsecase.GetUserByUsername(username)
	user.Token = token

	return repository.UpdateUser(user)
}

func (userUsecase *UserUsecase) GetAuth(jwtReceived string, httpMethodRequested string) (doAction bool, msg string) {

	doAction = true
	msg = "Success"

	if !checker.HasFieldsRequired([]string {jwtReceived}){
		msg = "EmptyJwtError"
		doAction = false
		return
	}

	fmt.Println("Estoy aquí (línea 156 Usecase.go) y este es el jwt: " + jwtReceived)

	err := userUsecase.FindJwt(jwtReceived)

	if err != nil {
		msg = "InvalidJwtError"
		doAction = false
		return
	}

	userRequesting, _ := userUsecase.GetUserByJwt(jwtReceived)

	if !checker.HasPermissions(userRequesting.Role, httpMethodRequested){
		msg = "UserPermissionsError"
		doAction = false
		return
	}

	return
}
