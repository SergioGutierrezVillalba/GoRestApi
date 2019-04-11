package entity

import (
	"FirstProject/Model"
)

type RepositoryInterface interface{
	GetAll() ([]model.Timer, error)
	GetById(id string) (model.Timer, error)
	GetAllByUserId(userId string)([] model.Timer, error)
	Create(model.Timer) error
	Update(model.Timer) error
	Delete(id string) error
	InsertStartTime(model.Timer) error
}