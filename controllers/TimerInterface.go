package controllers

import (
	"time"
	"net/http"
	"encoding/json"

	"FirstProject/core/usecases"
	"FirstProject/model"

	"gopkg.in/mgo.v2/bson"

	// "fmt"
	"strconv"
)

var (
	timerUsecase usecases.TimerUsecase
)

type TimerInterface struct{}


func (timerInterface *TimerInterface) GetTimers(response http.ResponseWriter, request *http.Request){

	timers, err := timerUsecase.GetTimers()

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	timersFormatted, err := formatTimersForResponse(timers)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(response, http.StatusOK, timersFormatted)
}

func (timerInterface *TimerInterface) InitCountTime(response http.ResponseWriter, request *http.Request){

	var user model.User
	var timer model.Timer

	json.NewDecoder(request.Body).Decode(&user)

	startTime := time.Now().Unix()

	timer.Id = bson.NewObjectId()
	timer.UserId = user.Id.Hex()
	timer.Start = startTime

	err := timerUsecase.InsertStartTime(timer)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(response, http.StatusOK, timer.Id)
}

func (timerInterface *TimerInterface) FinishCountTime(response http.ResponseWriter, request *http.Request){

	var timer model.Timer

	json.NewDecoder(request.Body).Decode(&timer)

	timerDb, err := timerUsecase.GetTimeById(timer.Id.Hex())

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	finishTime := time.Now().Unix()
	duration := finishTime - timerDb.Start

	timerDb.Finish = finishTime
	timerDb.Duration = duration

	err = timerUsecase.InsertFinishTime(timerDb)

	if err != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	timerFormatted, err2 := formatTimerForResponse(timerDb)

	if err2 != nil {
		respond.WithError(response, http.StatusBadRequest, err.Error())
		return
	}

	respond.WithJson(response, http.StatusOK, timerFormatted)
}



// Functionalities

func formatTimersForResponse(timers [] model.Timer) (timersFormatted [] model.TimerFormatted, err error){

	for i := 0; i < len(timers); i++ {
		timerFormatted, errFormat := formatTimerForResponse(timers[i])

		if errFormat != nil {
			err = errFormat
			return
		}
		
		timersFormatted = append(timersFormatted, timerFormatted)
	}

	return
}

func formatTimerForResponse(timer model.Timer) (timerFormatted model.TimerFormatted, err error){

	timerFormatted.Id = timer.Id
	timerFormatted.UserId = timer.UserId
	timerFormatted.Duration = timer.Duration
	timerFormatted.Start, err = unixDateToString(timer.Start)

	if err != nil {
		return timerFormatted, err
	}

	timerFormatted.Finish, err = unixDateToString(timer.Finish)

	if err != nil {
		return timerFormatted, err
	}

	return timerFormatted, err
}

func unixDateToString(miliseconds int64) (string, error) {

	miliToString := strconv.FormatInt(miliseconds, 10) // int64 to String
	i, err := strconv.ParseInt(miliToString, 10, 64) // String to int

    if err != nil {
        return "", err
	}
	tm := time.Unix(i, 0) // int to Unix timestamp
	tmToString := tm.String()

	return tmToString, nil
}