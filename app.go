package main

import (
	"fmt"
	"log"
	"net/http"
	"FirstProject/api/userapi"
	//   "encoding/json"
	//   "io/ioutil"

	//   "gopkg.in/mgo.v2/bson"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	var userapi userapi.UserAPI
	r := mux.NewRouter()
	r.HandleFunc("/users", userapi.FindAll).Methods("GET")
	r.HandleFunc("/user/{id}", userapi.Find).Methods("GET")
	r.HandleFunc("/users", userapi.Create).Methods("POST")
	r.HandleFunc("/login", userapi.Login).Methods("POST")
	r.HandleFunc("/user/delete/{id}", userapi.Delete).Methods("DELETE")
	r.HandleFunc("/users", userapi.Update).Methods("PUT")
	
	corsObj := handlers.AllowedOrigins([]string{"*"})

	if err := http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}), corsObj)(r)); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Listening on port 3000....")
	}
}
