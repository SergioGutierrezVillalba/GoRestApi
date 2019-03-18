package usecases

import (
	"FirstProject/core/entities"
	"FirstProject/model"
)

type UserUsecase struct {}

var repository entities.Repository


func (userUsecase *UserUsecase) GetUsers() ([]model.User, error){
	return repository.GetUsers()
}

func (userUsecase *UserUsecase) GetUserById(userId string) (interface{}, error){
	return repository.GetUserById(userId)
}

func (userUsecase *UserUsecase) GetUserByUsername(username string) (model.User, error){
	return repository.GetUserByUsername(username)
}

func (userUsecase *UserUsecase) GetUserByJwt(jwt string)(model.User, error){
	return repository.GetUserByJwt(jwt)
}

func (userUsecase *UserUsecase) GetEmailOfUser(username string) (model.User, error){
	return repository.GetEmailOfUser(username)
}

func (userUsecase *UserUsecase) GetUserByToken(token string) (model.User, error){
	return repository.GetUserByToken(token)
}

func (userUsecase *UserUsecase) FindJwt(jwtReceived string) error{
	return repository.FindJwt(jwtReceived)
}

func (userUsecase *UserUsecase) CreateUser(user model.User) error{
	return repository.CreateUser(user)
}

func (userUsecase *UserUsecase) UpdateUser(user model.User) error{
	return repository.UpdateUser(user)
}

func (userUsecase *UserUsecase) DeleteUser(userId string) error{
	return repository.DeleteUser(userId)
}

func (userUsecase *UserUsecase) SaveToken(token string, username string) error {

	user, _ := userUsecase.GetUserByUsername(username) // ya he comprobado antes si existia
	user.Token = token
	err := userUsecase.UpdateUser(user)

	return err

}

