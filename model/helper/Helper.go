package helper

import (
	"net/http"
	"strings"
	"FirstProject/model/auth"
	"FirstProject/model"
	// "fmt"
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
