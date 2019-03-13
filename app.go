package main

import (
	"fmt"
	"log"
	"net/http"

	"FirstProject/api/userapi"
	"FirstProject/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	config.InitializeConnection()

	var userapi userapi.UserAPI
	r := mux.NewRouter()

	// ROUTES
	r.HandleFunc("/users", userapi.FindAll).Methods("GET")
	r.HandleFunc("/users/{id}", userapi.Find).Methods("GET")

	r.HandleFunc("/users", userapi.Create).Methods("POST")
	r.HandleFunc("/login", userapi.Login).Methods("POST")
	r.HandleFunc("/sendrecover", userapi.SendRecover).Methods("POST")

	r.HandleFunc("/reset", userapi.Reset).Methods("PATCH")

	r.HandleFunc("/users", userapi.Update).Methods("PUT")

	r.HandleFunc("/users/{id}", userapi.Delete).Methods("DELETE")
	
	
	
	// BACKEND CONNECTION PERMISSIONS
	corsObj := handlers.AllowedOrigins([]string{"*"})

	if err := http.ListenAndServe(":3001", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}), corsObj)(r)); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Listening on port 3001....")
	}
}
