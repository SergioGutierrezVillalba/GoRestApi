package main
 
import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "io/ioutil"

  "FirstProject/model"

  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/mux"

)

var repo = model.Repository{Server:"http://localhost:27017", Database:"project"}


func AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
	repo.GetAll()
}

func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close();
	var user model.User
	// user := user.User{UserName: "yo"}

	jsn, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal("Error reading the body", err)
	}

	err = json.Unmarshal(jsn, &user) // put info into the user struct
	user.ID = bson.NewObjectId()

	
	if err:= repo.InsertUser(user); err != nil {
		log.Fatal("Error adding new register in database", err)
	} 
	

	// fmt.Fprintln(w, user);
	
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
	r.HandleFunc("/users/create", CreateUserEndPoint).Methods("POST")
	r.HandleFunc("/users", UpdateUserEndPoint).Methods("PUT")
	r.HandleFunc("/users", DeleteUserEndPoint).Methods("DELETE")
	r.HandleFunc("/users/{id}", FindUserEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}