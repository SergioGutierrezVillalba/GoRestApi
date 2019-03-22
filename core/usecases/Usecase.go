package usecases

import (
	"FirstProject/core/entities"
	"FirstProject/model/validation"
)

var (
	repository 		entities.Repository
	
	checker 		validation.Checker
)

type Usecase struct{}

func (usecase *Usecase) FindJwt(jwt string) error {
	return repository.FindJwt(jwt)
}