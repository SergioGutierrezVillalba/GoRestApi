package auth

import (
	"net/http"
	"fmt"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"FirstProject/Model"
	"FirstProject/Model/helper"
	"FirstProject/Model/auth"
	"FirstProject/Domains/user/usecase"
	repo "FirstProject/Domains/user/entity"

)

type GAuthToken struct {}

var (
	respond				model.Responser

	Helper				helper.Helper

	Auth				auth.Authentication

	usersRepo 			repo.RepositoryInterface
	usersUsecase		usecase.Usecase
)

func (gAuthToken *GAuthToken) Middleware(h http.Handler, session *mgo.Session, methodRequested string) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

		newSession := session.Copy()
		defer newSession.Close()

		SetConnectors(newSession)
		jwt := Helper.GetJWTFromHeader(r)

		if Helper.IsEmpty(jwt) {
			fmt.Println("(Middleware blocks): Empty JWT")
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		}

		uncryptedJWT := Auth.Decrypt(Auth.DecodeBase64(jwt))
		expiration := Auth.GetExpirationTimeOfJWT(uncryptedJWT)

		if Auth.IsNotValid(expiration){
			fmt.Println("(Middleware says): Expirated time: " + string(expiration))
			fmt.Println("(Middleware blocks): Expirated JWT")
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		}

		userRequesting, err := usersUsecase.GetUserByJwt(jwt)
		fmt.Println("Rol:" + userRequesting.Role)
		
		if err != nil {
			fmt.Println("(Middleware blocks): JWT not exists")
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		}

		r.Header.Add("User-Agent", userRequesting.Role)

		if !HasPermissionForDoThatRequest(userRequesting.Role, methodRequested){
			fmt.Println("(Middleware blocks): Not enough permissions")
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		}
		
		h.ServeHTTP(w, r)
	})
}

func SetConnectors(session *mgo.Session){
	usersRepo = repo.NewMongoDbRepository(session)
	usersUsecase = usecase.NewUsecase(usersRepo)
}

func HasPermissionForDoThatRequest(roleVerified string, methodRequested string) bool {

	if roleVerified == "user" {
		roleVerified = "self"
	}


	Permissions	:= map[string][]string{
		"GetAllUsers":{"ADMIN"},
		"GetUserById":{"ADMIN"},
		"GetMe":{"ADMIN", "SELF"},
		"CreateUser":{"ADMIN"},
		"UpdateUser":{"ADMIN", "SELF"},
		"UpdateUserWithoutPassword":{"ADMIN", "SELF"},
		"DeleteUser":{"ADMIN"},
		"SetProfileImage":{"ADMIN", "SELF"},
		"GetProfileImage":{"ADMIN", "SELF"},

		"GetAllTimers":{"ADMIN"},
		"GetTimerById":{"ADMIN"},
		"GetTimersByUserId":{"ADMIN", "SELF"}, // needs more verification in self
		"CreateTimer":{"ADMIN"},
		"UpdateTimer":{"ADMIN"},
		"DeleteTimer":{"ADMIN"},
		"StartTimer":{"ADMIN", "SELF"},
		"FinishTimer":{"ADMIN", "SELF"},
	}

	for _, roleAllowed := range Permissions[methodRequested]{
		if roleAllowed == strings.ToUpper(roleVerified) {
			return true
		}
	}
	return false
}
