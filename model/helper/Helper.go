package helper

import (
	"net/http"
	"strings"
	"strconv"
	"errors"
	"time"
	"fmt"

	"FirstProject/Model/Auth"
	"FirstProject/Model"	

	"github.com/gorilla/mux"
)

var(
	Auth 	auth.Authentication
)

type Helper struct{}

func (h *Helper) GetJWTFromHeaderRequest(r *http.Request) (cleanJWT string, err error){

	err = errors.New("HeaderError")
	rawJWT := r.Header.Get("Authorization")

	if h.IsNotEmpty(rawJWT) {
		if h.DoesntHasBearer(rawJWT) {
			return
		}
		cleanJWT = QuitBearer(rawJWT)
		err = nil
	}
	return
}

func (h *Helper) IsNotEmpty(data string) bool {
	if data == "" {
		return false
	}
	return true
}

func (h *Helper) DoesntHasBearer(jwt string) bool {
	return !strings.Contains(jwt, "Bearer")
}

// DO NOT QUIT whitespace before Bearer, is neccesary
// for removing whitespace after word 'Bearer'
func QuitBearer(rawJWT string)(cleanJWT string){
	cleanJWT = strings.Replace(rawJWT, "Bearer ", "" , 1)
	fmt.Println(cleanJWT)
	return
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

func (h *Helper) IsEmpty(data string) bool {
	if data == "" {
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

func (h *Helper) GetIdFromUrl(request *http.Request) (id string) {
	vars := mux.Vars(request)
	id = vars["id"]
	return
}

func (h *Helper) GetDateFromUrl(request *http.Request) (date int64) {
	vars := mux.Vars(request)
	dateStringed := vars["date"]
	
	n, err := strconv.ParseInt(dateStringed, 10, 64)

	if err != nil {
		return
	}

	date = n
	return
}

func (h *Helper) ActionGivesError(err error) bool {
	if err != nil {
		return true
	}
	return false
}

// We have to determine what user is Requesting and
// which is susceptible to change to know which Role
// has the user Requesting and in consequence what it
// is allowed to do
func (h *Helper) WhichRoleIsUsed (userRequesting model.User, userToModify model.User) (situation string) {
	if h.IsUser(userRequesting) {
		if !h.IsUpdatingItself(userRequesting.Id.Hex(), userToModify.Id.Hex()){
			fmt.Println("(WhichRoleIsUsed): You arent allowed to do that")
			return "NOAUTH"
		}
		return "SELF"
	}
	return "ADMIN"
}

func (h *Helper) CleanSlashesFromToken(token string) string {
	tokenSanitized := strings.Replace(token, "/", "", -1)
	return tokenSanitized
}