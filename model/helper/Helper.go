package helper

import (
	"net/http"
	"strings"
	"time"
	"fmt"

	"FirstProject/Model/auth"
	"FirstProject/Model"	
)

var(
	Auth 	auth.Authentication
)

type Helper struct{}

func (h *Helper) GetJWTFromHeader(r *http.Request) (cleanJWT string){

	cleanJWT = ""
	rawJWT:= r.Header["Authorization"][0]

	if h.IsNotEmpty(rawJWT) {
		if h.DoesntHasBearer(rawJWT) {
			return
		}
		withoutBearer := strings.Split(rawJWT, "Bearer")[1]
		cleanJWT = strings.Trim(withoutBearer, " ")
	}
	return
}

// NEEDED BECAUSE IT CAN MAKE PANIC IF NOT ARRIVES
// Need to fix, " before token causes panic too
func (h *Helper) DoesntHasBearer(jwt string) bool {
	return !strings.Contains(jwt, "Bearer")
}

func (h *Helper) IsUser(user model.User) bool {
	if user.Role == "user" {
		return true
	}
	return false
}

func (h *Helper) IsAdmin(user model.User) bool {
	if user.Role == "admin" {
		return true
	}
	return false
}

func (h *Helper) IsUpdatingItself(userIdRequesting string, userIdUpdating string) bool {

	if userIdRequesting == userIdUpdating {
		return true
	}
	return false
}

func (h *Helper) IsNotEmpty(data string) bool {
	if data == "" {
		return false
	}
	return true
}

func (h *Helper) IsEmpty(data string) bool {
	if data == "" {
		return true
	}
	return false
}

func (h *Helper) IsAlreadyFinished(timerFinish int64) bool {
	if timerFinish > 0 {
		return true
	}
	return false
}

func (h *Helper) UnixDateToString(miliseconds int64) (dateStringed string) {
	dateStringed = time.Unix(miliseconds, 0).String()
	return
}

func (h *Helper) FormatTimerForResponse(timer model.Timer) (timerFormatted model.TimerFormatted){

	timerFormatted.Id = timer.Id
	timerFormatted.UserId = timer.UserId
	timerFormatted.Duration = timer.Duration
	timerFormatted.Start = h.UnixDateToString(timer.Start)
	timerFormatted.Finish = h.UnixDateToString(timer.Finish)
	return
}

func (h *Helper) FormatTimersForResponse(timers [] model.Timer) (timersFormatted [] model.TimerFormatted){
	for i := 0; i < len(timers); i++ {
		timerFormatted := h.FormatTimerForResponse(timers[i])
		timersFormatted = append(timersFormatted, timerFormatted)
	}
	return
}

// We have to determine what user is Requesting and
// which is susceptible to change to know which Role
// has the user Requesting and in consequence what it
// is allowed to do

func (h *Helper) WhichRoleIsUsed (userRequesting model.User, userToModify model.User) (situation string) {
	if h.IsUser(userRequesting) {
		if !h.IsUpdatingItself(userRequesting.Id.Hex(), userToModify.Id.Hex()){
			fmt.Println("(SelfRoleRequest): You arent allowed to do that")
			return "NOAUTH"
		}
		return "SELF"
	}
	return "ADMIN"
}