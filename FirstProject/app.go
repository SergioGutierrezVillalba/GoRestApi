package main
 
import (
  "fmt"
  "log"
  "net/http"
  "model/user"

  "github.com/gorilla/mux"
  // "gopkg.in/mgo.v2"
  // "gopkg.in/mgo.v2/bson"
)


func AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {
  user := user.User{}
	fmt.Fprintln(w, "not implemented yet !")
}

func UpdateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func DeleteUserEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", AllUsersEndPoint).Methods("GET")
	r.HandleFunc("/users", CreateUserEndPoint).Methods("POST")
	r.HandleFunc("/users", UpdateUserEndPoint).Methods("PUT")
	r.HandleFunc("/users", DeleteUserEndPoint).Methods("DELETE")
	r.HandleFunc("/users/{id}", FindUserEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}