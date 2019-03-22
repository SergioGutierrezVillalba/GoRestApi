package main

import (
	"fmt"
	"log"
	"net/http"

	"FirstProject/controllers"
	"FirstProject/model/database"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

var (
	userInterface 	controllers.UserInterface
	timerInterface 	controllers.TimerInterface
)

func main(){
	StartDBConnection()
	StartApi()
}

func StartDBConnection(){
	database.StartConnection()
}

func StartApi(){

	r := mux.NewRouter()

	SetRoutesToApi(r)
	StartServerAndAllowConections(r)
}

func SetRoutesToApi(r *mux.Router){
	// USERS
	r.HandleFunc("/users", userInterface.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userInterface.GetUserById).Methods("GET")
	r.HandleFunc("/users", userInterface.CreateUser).Methods("POST")
	r.HandleFunc("/login", userInterface.Login).Methods("POST")
	r.HandleFunc("/sendrecover", userInterface.SendRecover).Methods("POST")
	r.HandleFunc("/reset", userInterface.Reset).Methods("PATCH")
	r.HandleFunc("/users", userInterface.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userInterface.DeleteUser).Methods("DELETE")


	// TIMER
	r.HandleFunc("/timers", timerInterface.GetTimers).Methods("GET")
	// r.HandleFunc("/timers/{id}", timerInterface.GetTimerById).Methods("GET")
	// r.HandleFunc("/timers/users/{id}", timerInterface.GetTimersByUserId).Methods("POST")
	r.HandleFunc("/startTimer", timerInterface.InitCountTime).Methods("POST")
	r.HandleFunc("/finishTimer", timerInterface.FinishCountTime).Methods("POST")
}

func StartServerAndAllowConections(r *mux.Router){

	corsObj := handlers.AllowedOrigins([]string{"*"})

	if err := http.ListenAndServe(":3003", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}), corsObj)(r)); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Listening on port 3003....")
	}
}