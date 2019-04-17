package entity

import (
	"FirstProject/Model"
)

type RepositoryInterface interface{
	GetAll() ([]model.User, error)
	GetById(string) (model.User, error)
	GetUserByJwt(string)(model.User, error)
	GetUserByUsername(string)(model.User, error)
	GetUserByRecoverToken(string)(model.User, error)
	Create(model.User) error
	Update(model.User) error
	Delete(string) error
}