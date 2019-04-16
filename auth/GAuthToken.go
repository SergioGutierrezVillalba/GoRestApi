package auth

import (
	"net/http"
	"fmt"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"FirstProject/Model"
	"FirstProject/Model/helper"
	"FirstProject/Model/Auth"
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

		StartRepositoriesAndUsecases(newSession /*repos slice*/)

		jwt := Helper.GetJWTFromHeader(r)

		if VerificationIsDenied(w, r, jwt, methodRequested) {
			respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		}
		
		h.ServeHTTP(w, r)
	})
}

// usersRepo must be changed for a slice of repos
func StartRepositoriesAndUsecases(session *mgo.Session/* repos slice */){
	SetSessionToRepositories(session)
	SetReposToUsecases(/* repos slice */)
}
func SetSessionToRepositories(session *mgo.Session){
	usersRepo = repo.NewMongoDbRepository(session)
}
func SetReposToUsecases(/*repos*/){
	// for --> iterate over repos slice and create a map with ["context"] = newUsecase
	// but if i only have to use users usecase for verification?
	usersUsecase = usecase.NewUsecase(usersRepo)
}
func VerificationIsDenied(w http.ResponseWriter, r *http.Request, jwt string, methodRequested string) bool {

	if UserIsNotAllowed(jwt, methodRequested){
		return true
	}
	if JWTIsNotValid(w, r, jwt) {
		return true
	}
	return false 
}
func JWTIsNotValid(w http.ResponseWriter, r *http.Request, jwt string)bool{

	uncryptedJWT := Auth.Decrypt(Auth.DecodeBase64(jwt))
	expiration := Auth.GetExpirationTimeOfJWT(uncryptedJWT)

	if Auth.IsExpirated(expiration){
		fmt.Println("(Middleware blocks): Expirated JWT")
		return true
	}
	return false
}
func UserIsNotAllowed(jwt string, methodRequested string) bool {

	userRequesting, _ := usersUsecase.GetUserByJwt(jwt)
	fmt.Println("Rol:" + userRequesting.Role)
		
	if userRequesting.NotExists() {
		fmt.Println("(Middleware blocks): User not exists")
		return true
	}

	if !HasPermissionForDoThatRequest(userRequesting.Role, methodRequested){ // Doesnt has permissions
		fmt.Println("(Middleware blocks): Not enough permissions")
		return true
	}
	return false
}
func HasPermissionForDoThatRequest(roleVerified string, methodRequested string) bool {

	if roleVerified == "user" { // must be changed
		roleVerified = "self"
	}

	Permissions	:= map[string][]string{ // to an external file
		"GetAllUsers":{"ADMIN"},
		"GetUserById":{"ADMIN"},
		"GetMe":{"ADMIN", "SELF"},
		"CreateUser":{"ADMIN"},
		"UpdateUser":{"ADMIN", "SELF"},
		"UpdateUserWithoutUpdatingPassword":{"ADMIN", "SELF"},
		"DeleteUser":{"ADMIN"},
		"SetProfileImage":{"ADMIN", "SELF"},
		"GetProfileImage":{"ADMIN", "SELF"},

		"GetAllTimers":{"ADMIN"},
		"GetTimerById":{"ADMIN"},
		"GetTimersByUserId":{"ADMIN", "SELF"},
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
