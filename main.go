package main

import (
	"fmt"
	"log"
	"net/http"
	
	"FirstProject/model/database"
	"FirstProject/auth"

	usersController "FirstProject/Domains/user/controller"

	usersUsecase "FirstProject/Domains/user/usecase"
	timersUsecase "FirstProject/Domains/timer/usecase"

	usersRepo "FirstProject/Domains/user/entity"
	timersRepo "FirstProject/Domains/timer/entity"
	
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

var (
	db				database.Db
	gAuthToken		auth.GAuthToken
)

func main(){
	db.StartConnection()
	StartApi()
}

func StartApi(){

	r := mux.NewRouter()

	// REPOS
	usersRepo := usersRepo.NewMongoDbRepository(db.Session)
	timersRepo := timersRepo.NewMongoDbRepository(db.Session)

	// USECASES
	usersUsecase := usersUsecase.NewUsecase(usersRepo)
	timersUsecase := timersUsecase.NewUsecase(timersRepo)

	// INTERFACES
	usersController := usersController.NewUsersController(usersUsecase, timersUsecase)


	// USERS HANDLERS
	getMe := gAuthToken.Middleware(http.HandlerFunc(usersController.GetUserByJwt), db.Session, "GetMe")
	getUsers := gAuthToken.Middleware(http.HandlerFunc(usersController.GetAllUsers), db.Session, "GetAllUsers")
	getUserById := gAuthToken.Middleware(http.HandlerFunc(usersController.GetUserById), db.Session, "GetUserById")
	createUser := gAuthToken.Middleware(http.HandlerFunc(usersController.CreateUser), db.Session, "CreateUser")
	updateUser := gAuthToken.Middleware(http.HandlerFunc(usersController.UpdateUser), db.Session, "UpdateUser")
	deleteUser := gAuthToken.Middleware(http.HandlerFunc(usersController.DeleteUser), db.Session, "DeleteUser")
	register := http.HandlerFunc(usersController.Register)
	login := http.HandlerFunc(usersController.Login)
	sendRecover := http.HandlerFunc(usersController.SendRecover)
	resetPassword := http.HandlerFunc(usersController.ResetPassword)

	// TIMERS HANDLERS
	getTimers := gAuthToken.Middleware(http.HandlerFunc(usersController.GetAllTimers), db.Session, "GetAllTimers")
	getTimerById := gAuthToken.Middleware(http.HandlerFunc(usersController.GetTimerById), db.Session, "GetTimerById")
	getTimersByUserId := gAuthToken.Middleware(http.HandlerFunc(usersController.GetTimersByUserId), db.Session, "GetTimersByUserId")
	createTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.CreateTimer), db.Session, "CreateTimer")
	updateTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.UpdateTimer), db.Session, "UpdateTimer")
	deleteTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.DeleteTimer), db.Session, "DeleteTimer")
	startTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.StartTimer), db.Session, "StartTimer")
	finishTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.FinishTimer), db.Session, "FinishTimer")

	// USERS ROUTES
	r.Handle("/user", getMe).Methods("GET")
	r.Handle("/users", getUsers).Methods("GET")
	r.Handle("/users/{id}", getUserById).Methods("GET")
	r.Handle("/users", createUser).Methods("POST")
	r.Handle("/users", updateUser).Methods("PUT")
	r.Handle("/users/{id}", deleteUser).Methods("DELETE")
	r.Handle("/users/login", login).Methods("POST")
	r.Handle("/users/register", register).Methods("POST")
	r.Handle("/users/sendrecover", sendRecover).Methods("POST")
	r.Handle("/users/reset", resetPassword).Methods("PATCH")


	// TIMERS ROUTES
	r.Handle("/timers", getTimers).Methods("GET")
	r.Handle("/timers/{id}", getTimerById).Methods("GET")
	r.Handle("/timers/users/{id}", getTimersByUserId).Methods("GET")
	r.Handle("/timers", createTimer).Methods("POST")
	r.Handle("/timers", updateTimer).Methods("PUT")
	r.Handle("/timers/{id}", deleteTimer).Methods("DELETE")
	r.Handle("/timers/start", startTimer).Methods("POST")
	r.Handle("/timers/finish", finishTimer).Methods("POST")

	StartServerAndAllowConections(r)

}

func StartServerAndAllowConections(r *mux.Router){

	corsObj := handlers.AllowedOrigins([]string{"*"})

	if err := http.ListenAndServe(":3003", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}), corsObj)(r)); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Listening on port 3003....")
	}
}