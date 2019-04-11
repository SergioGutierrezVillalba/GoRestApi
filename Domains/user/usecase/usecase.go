package usecase

import (
	repo "FirstProject/Domains/user/entity"
	"FirstProject/model"
	"FirstProject/model/validation"

	"errors"
	// "fmt"
)

type Usecase interface{
	GetAll() ([] model.User, error)
	GetById(string) (model.User, error)
	GetUserByJwt(string)(model.User, error)
	GetUserByUsername(string)(model.User, error)
	GetUserByRecoverToken(string)(model.User, error)
	Create(model.User) error
	Update(model.User) error
	UpdateSelf(model.User) error
	UpdateAdmin(model.User) error
	Delete(string) error
	Register(model.User) error
}

var (
	checker 	validation.Checker
)

type UsersUsecase struct {
	repo 	repo.RepositoryInterface
}

func NewUsecase(r repo.RepositoryInterface) Usecase {
	return &UsersUsecase{
		repo: r,
	}
}

func (u *UsersUsecase) GetAll()(users [] model.User, err error){

	users, err = u.repo.GetAll()

	if err != nil {
		err = errors.New("UsersNotFoundError")
	}

	return
}

func (u *UsersUsecase) GetById(userId string)(user model.User, err error){

	user, err = u.repo.GetById(userId)

	if err != nil {
		err = errors.New("UserNotFoundError")
	}

	return
}

func (u *UsersUsecase) GetUserByUsername(username string)(user model.User, err error){
	user, err = u.repo.GetUserByUsername(username)

	if err != nil {
		err = errors.New("UserNotFoundError")
	}
	return
}

func (u *UsersUsecase) GetUserByJwt(jwt string)(user model.User, err error){
	user, err = u.repo.GetUserByJwt(jwt)

	if err != nil {
		err = errors.New("UserNotFoundError")
	}
	return
}

func (u *UsersUsecase) GetUserByRecoverToken(token string)(user model.User, err error){
	user, err = u.repo.GetUserByRecoverToken(token)

	if err != nil {
		err = errors.New("UserNotFoundError")
	}
	return
}

func (u *UsersUsecase) Create(user model.User) (err error) {

	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email, user.Role)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	_, err = u.GetUserByUsername(user.Username)

	if err == nil {
		err = errors.New("UsernameAlreadyExistsError")
		return
	}

	err = u.repo.Create(user)

	if err != nil {
		err = errors.New("CreationUserError")
		return
	}
	return
}

func (u *UsersUsecase) Update(user model.User) (err error) {

	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	_, err = u.repo.GetById(user.Id.Hex())

	if err != nil {
		err = errors.New("UserNotFoundError")
		return
	}

	err = u.repo.Update(user)

	if err != nil {
		err = errors.New("UpdateUserError")
		return
	}
	return
}

func (u *UsersUsecase) UpdateSelf(user model.User) (err error){

	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email, user.GroupId)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	if err = u.repo.Update(user); err != nil {
		err = errors.New("UpdateUserError")
		return
	}
	return
}

func (u *UsersUsecase) UpdateAdmin(user model.User) (err error){

	var fieldsRequired  []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Role, user.Email, user.GroupId)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	if err = u.repo.Update(user); err != nil {
		err = errors.New("UpdateUserError")
		return
	}
	return 
}

func (u *UsersUsecase) Delete(userId string) (err error) {
	err = u.repo.Delete(userId)

	if err != nil {
		err = errors.New("DeleteUserError")
	}
	return
}

func (u *UsersUsecase) Register(user model.User) (err error){
	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, user.Username, user.Password, user.Email)

	if !checker.HasFieldsRequired(fieldsRequired){
		err = errors.New("EmptyFieldsError")
		return
	}

	_, err = u.GetUserByUsername(user.Username)

	if err == nil {
		err = errors.New("UsernameAlreadyExistsError")
		return
	}

	user.Role = "user"
	err = u.repo.Create(user)

	if err != nil {
		err = errors.New("CreationUserError")
		return
	}
	return
}