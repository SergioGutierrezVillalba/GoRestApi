package main

import (
	"fmt"
	"log"
	"net/http"

	"FirstProject/presenters"
	"FirstProject/model/database"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func main(){

	StartDBConnection()
	StartApi()

}

func StartDBConnection(){
	database.StartConnection()
}

func StartApi(){

	var userInterface presenters.UserInterface

	r := mux.NewRouter()

	// ROUTES
	r.HandleFunc("/users", userInterface.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userInterface.GetUserById).Methods("GET")

	r.HandleFunc("/users", userInterface.CreateUser).Methods("POST")
	r.HandleFunc("/login", userInterface.Login).Methods("POST")
	r.HandleFunc("/sendrecover", userInterface.SendRecover).Methods("POST")

	r.HandleFunc("/reset", userInterface.Reset).Methods("PATCH")

	r.HandleFunc("/users", userInterface.UpdateUser).Methods("PUT")

	r.HandleFunc("/users/{id}", userInterface.DeleteUser).Methods("DELETE")
	
	// BACKEND CONNECTION PERMISSIONS
	corsObj := handlers.AllowedOrigins([]string{"*"})

	if err := http.ListenAndServe(":3002", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}), corsObj)(r)); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Listening on port 3002....")
	}

}