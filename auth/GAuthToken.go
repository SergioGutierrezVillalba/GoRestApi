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
	Respond				model.Responser

	Helper				helper.Helper

	Auth				auth.Authentication

	UsersRepo 			repo.RepositoryInterface
	UsersUsecase		usecase.Usecase

	NewSession			*mgo.Session

	ResponseWriter		http.ResponseWriter
	Request				*http.Request		
	
	JWTUsed				string
)

func (gAuthToken *GAuthToken) Middleware(h http.Handler, session *mgo.Session, methodRequested string) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

		if ActionGivesError(SetGlobalVars(w, r, session)) {
			Respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		}

		InitRepoAndUsecaseRoleCheckers()
		defer NewSession.Close()

		if VerificationIsDeniedFor(methodRequested) {
			Respond.WithError(w, http.StatusBadRequest, "Unauthorized")
			return
		}
		
		h.ServeHTTP(w, r)
	})
}

func ActionGivesError(err error) bool {
	return Helper.ActionGivesError(err)
}

func SetGlobalVars(w http.ResponseWriter, r *http.Request, session *mgo.Session)(err error){
	ResponseWriter = w
	Request = r
	NewSession = session.Copy()
	JWTUsed, err = Helper.GetJWTFromHeaderRequest(r)
	return
}

func InitRepoAndUsecaseRoleCheckers(){
	SetSessionToRepo()
	SetRepoToUsecase()
}

func SetSessionToRepo(){
	UsersRepo = repo.NewMongoDbRepository(NewSession)
}

func SetRepoToUsecase(){
	UsersUsecase = usecase.NewUsecase(UsersRepo)
}

func VerificationIsDeniedFor(methodRequested string) bool {

	if UserIsNotAllowed(methodRequested){
		return true
	}
	if JWTIsNotValid() {
		return true
	}
	return false 
}

func UserIsNotAllowed(methodRequested string) bool {

	userRequesting, _ := UsersUsecase.GetUserByJwt(JWTUsed)
	fmt.Println("Rol:" + userRequesting.Role)
		
	if userRequesting.NotExists() {
		fmt.Println("(Middleware blocks): User not exists")
		return true
	}

	// TODO DoesntHasPermissions()
	if !HasPermissionForDoThatRequest(userRequesting.Role, methodRequested){ 
		fmt.Println("(Middleware blocks): Not enough permissions")
		return true
	}
	return false
}

func JWTIsNotValid() bool {

	uncryptedJWT := Auth.Decrypt(Auth.DecodeBase64(JWTUsed))
	expiration := Auth.GetExpirationTimeOfJWT(uncryptedJWT)

	if Auth.IsExpirated(expiration){
		fmt.Println("(Middleware blocks): Expirated JWT")
		return true
	}
	return false
}

func HasPermissionForDoThatRequest(roleVerified string, methodRequested string) bool {

	// TODO must be changed when self and admin at the same time, have a problem
	if roleVerified == "user" { 
		roleVerified = "self"
	}

	// TODO to an external file
	Permissions	:= map[string][]string{ 
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

	// TODO CheckPermissionInPermissionsList()
	for _, roleAllowed := range Permissions[methodRequested]{
		if roleAllowed == strings.ToUpper(roleVerified) {
			return true
		}
	}
	return false
}
