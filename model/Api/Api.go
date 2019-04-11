package api

import (
	"log"
	"fmt"
	"net/http"

	auth "FirstProject/auth"

	usersController "FirstProject/Domains/user/controller"

	usersUsecase "FirstProject/Domains/user/usecase"
	timersUsecase "FirstProject/Domains/timer/usecase"

	usersRepo "FirstProject/Domains/user/entity"
	timersRepo "FirstProject/Domains/timer/entity"

	mgo "gopkg.in/mgo.v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

type Api struct {}

var (
	gAuthToken		auth.GAuthToken
)

func (a *Api) Start(session *mgo.Session){

	r := mux.NewRouter()

	// REPOS
	usersRepo := usersRepo.NewMongoDbRepository(session)
	timersRepo := timersRepo.NewMongoDbRepository(session)


	// USECASES
	usersUsecase := usersUsecase.NewUsecase(usersRepo)
	timersUsecase := timersUsecase.NewUsecase(timersRepo)


	// INTERFACES
	usersController := usersController.NewUsersController(usersUsecase, timersUsecase)


	// USERS HANDLERS
	getMe := gAuthToken.Middleware(http.HandlerFunc(usersController.GetUserByJwt), session, "GetMe")
	getUsers := gAuthToken.Middleware(http.HandlerFunc(usersController.GetAllUsers), session, "GetAllUsers")
	getUserById := gAuthToken.Middleware(http.HandlerFunc(usersController.GetUserById), session, "GetUserById")
	createUser := gAuthToken.Middleware(http.HandlerFunc(usersController.CreateUser), session, "CreateUser")
	updateUser := gAuthToken.Middleware(http.HandlerFunc(usersController.UpdateUser), session, "UpdateUser")
	updateUserWithoutPassword := gAuthToken.Middleware(http.HandlerFunc(usersController.UpdateUserWithoutPassword), session, "UpdateUserWithoutPassword")
	deleteUser := gAuthToken.Middleware(http.HandlerFunc(usersController.DeleteUser), session, "DeleteUser")
	// setProfileImageToUser := gAuthToken.Middleware(http.HandlerFunc(usersController.SetProfileImage), session, "SetProfileImage")
	register := http.HandlerFunc(usersController.Register)
	login := http.HandlerFunc(usersController.Login)
	sendRecover := http.HandlerFunc(usersController.SendRecover)
	resetPassword := http.HandlerFunc(usersController.ResetPassword)


	// TIMERS HANDLERS
	getTimers := gAuthToken.Middleware(http.HandlerFunc(usersController.GetAllTimers), session, "GetAllTimers")
	getTimerById := gAuthToken.Middleware(http.HandlerFunc(usersController.GetTimerById), session, "GetTimerById")
	getTimersByUserId := gAuthToken.Middleware(http.HandlerFunc(usersController.GetTimersByUserId), session, "GetTimersByUserId")
	createTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.CreateTimer), session, "CreateTimer")
	updateTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.UpdateTimer), session, "UpdateTimer")
	deleteTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.DeleteTimer), session, "DeleteTimer")
	startTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.StartTimer), session, "StartTimer")
	finishTimer := gAuthToken.Middleware(http.HandlerFunc(usersController.FinishTimer), session, "FinishTimer")


	// WEBSOCKET HANDLERS
	r.Handle("/", http.FileServer(http.Dir("./front")))
	createWebsocket := http.HandlerFunc(usersController.StartWebSocket)
	finishWebsocket := http.HandlerFunc(usersController.FinishWebSocket)
	

	// USERS ROUTES
	r.Handle("/user", getMe).Methods("GET")
	r.Handle("/users", getUsers).Methods("GET")
	r.Handle("/users/{id}", getUserById).Methods("GET")
	r.Handle("/users", createUser).Methods("POST")
	r.Handle("/users", updateUser).Methods("PUT")
	r.Handle("/users/nopwd", updateUserWithoutPassword).Methods("PUT")
	r.Handle("/users/{id}", deleteUser).Methods("DELETE")
	r.Handle("/users/login", login).Methods("POST")
	r.Handle("/users/register", register).Methods("POST")
	r.Handle("/users/sendrecover", sendRecover).Methods("POST")
	r.Handle("/users/reset", resetPassword).Methods("PATCH")
	
	// r.Handle("/")


	// TIMERS ROUTES
	r.Handle("/timers", getTimers).Methods("GET")
	r.Handle("/timers/{id}", getTimerById).Methods("GET")
	r.Handle("/timers/users/{id}", getTimersByUserId).Methods("GET")
	r.Handle("/timers", createTimer).Methods("POST")
	r.Handle("/timers", updateTimer).Methods("PUT")
	r.Handle("/timers/{id}", deleteTimer).Methods("DELETE")
	r.Handle("/timers/start", startTimer).Methods("POST")
	r.Handle("/timers/finish", finishTimer).Methods("POST")


	// WEBSOCKET ROUTES
	r.HandleFunc("/ws", createWebsocket)
	r.HandleFunc("/wsclose", finishWebsocket)

	StartServerAndAllowConections(r)
}

func StartServerAndAllowConections(r *mux.Router){

	corsObj := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})

	err := http.ListenAndServe(":3003", handlers.CORS( allowedHeaders, allowedMethods, corsObj )(r))

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Listening on port 3003....")
	}
}