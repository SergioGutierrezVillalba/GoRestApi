package usecase

import (
	repo "FirstProject/Domains/timer/entity"
	"FirstProject/Model"
	"FirstProject/Model/Validation"

	"errors"
	// "log"
)

var (
	Checker		validation.Checker
)

type Usecase interface{
	GetAll()([] model.Timer, error)
	GetById(string)(model.Timer, error)
	GetAllByUserId(string)([] model.Timer, error)
	Create(model.Timer)error
	Update(model.Timer)error
	Delete(string)error
	StartTimer(timer model.Timer)error
	FinishTimer(timer model.Timer) error
}

type TimerUsecase struct {
	repo 	repo.RepositoryInterface
}

func NewUsecase(r repo.RepositoryInterface) Usecase {
	return &TimerUsecase{
		repo: r,
	}
}

func (t *TimerUsecase) GetAll()(timers [] model.Timer, err error){
	timers, err = t.repo.GetAll()

	if err != nil {
		err = errors.New("TimersNotFoundError")
	}
	return
}
func (t *TimerUsecase) GetById(timerId string)(timer model.Timer, err error){
	timer, err = t.repo.GetById(timerId)

	if err != nil {
		err = errors.New("TimerNotFoundError")
	}
	return
}
func (t *TimerUsecase) GetAllByUserId(userId string)(timers [] model.Timer, err error){
	timers, err = t.repo.GetAllByUserId(userId)
	if err != nil {
		err = errors.New("TimersNotFound")
		return
	}
	return
}
func (t *TimerUsecase) Create(timer model.Timer) error {

	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, timer.UserId.Hex())

	if !Checker.HasFieldsRequired(fieldsRequired) {
		return errors.New("EmptyFieldsError")
	}

	err := t.repo.Create(timer)

	if err != nil {
		return errors.New("TimerCreationError")
	}
	return nil
}
func (t *TimerUsecase) Update(timer model.Timer) error {
	var fieldsRequired []string
	fieldsRequired = append(fieldsRequired, timer.UserId.Hex())

	if !Checker.HasFieldsRequired(fieldsRequired) {
		return errors.New("EmptyFieldsError")
	}

	timerDb, err := t.GetById(timer.GetId())

	if timerDb.NotExists() {
		return err
	}

	return t.repo.Update(timer)
}
func (t *TimerUsecase) Delete(timerId string)(err error){
	return t.repo.Delete(timerId)
}
func (t *TimerUsecase) StartTimer(timer model.Timer)(err error){
	err = t.repo.InsertStartTime(timer)
	return
}
func (t *TimerUsecase) FinishTimer(timer model.Timer) error {
	return t.repo.Update(timer)
}