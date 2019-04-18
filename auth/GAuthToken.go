package auth

import (
	"net/http"
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"

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

	if DoesntHasPermissionForDoThatRequest(userRequesting.Role, methodRequested){ 
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

func DoesntHasPermissionForDoThatRequest(roleVerified string, methodRequested string) bool {

	if roleVerified == "user" { 
		roleVerified = "self"
	}

	permissionsList := GetPermissionsList()
	hasPermission := CheckInPermissionsListIfHasPermissions(permissionsList, roleVerified, methodRequested)
	
	if hasPermission {
		return false
	}
	return true
}

func GetPermissionsList()(permissionsList map[string][]string){
	file, _ := ioutil.ReadFile("./permissions.json")
	err := json.Unmarshal([]byte(file), &permissionsList)
	if err != nil {
		fmt.Println("Error during permissions reading")
	}
	return
}

func CheckInPermissionsListIfHasPermissions(PermissisionsList map[string][]string, roleVerified string, methodRequested string) bool {
	for _, roleAllowed := range PermissisionsList[methodRequested]{
		if roleAllowed == strings.ToUpper(roleVerified) {
			return true
		}
	}
	return false
}

