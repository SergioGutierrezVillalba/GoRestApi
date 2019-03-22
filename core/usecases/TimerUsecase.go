package usecases

import (
	"errors"
	"FirstProject/model"
	// "fmt"
)



type TimerUsecase struct{}


func (timerUsecase *TimerUsecase) GetTimers() ([] model.Timer, error){
	return repository.GetTimers()
}

func (timerUsecase *TimerUsecase) GetTimeById(timeId string) (model.Timer, error) {

	var timer model.Timer

	if timeId == "" {
		return timer, errors.New("EmptyIdError")
	}

	return repository.GetTimeById(timeId)
}

func (timerUsecase *TimerUsecase) InsertStartTime(timer model.Timer) error{
	return repository.InsertStartTime(timer)
} 

func (timerUsecase *TimerUsecase) InsertFinishTime(timer model.Timer) error {

	timerDb, _ := timerUsecase.GetTimeById(timer.Id.Hex())

	if timerDb.Duration > 0 {
		return errors.New("TimeAlreadyFinished")
	}

	return repository.InsertFinishTime(timer)
}