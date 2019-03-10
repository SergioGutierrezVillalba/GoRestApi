package main
 
import (
  "fmt"
  "log"
  "net/http"
//   "encoding/json"
//   "io/ioutil"

//   "FirstProject/model"
  "FirstProject/api/userapi"

//   "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"

)

// var repo = model.Repository{Server:"mongodb://localhost:27017", Database:"project"}


// func AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
// 	repo.GetAll()
// }

// func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "not implemented yet !")
// }

// func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {

// 	defer r.Body.Close();
// 	var user model.User
// 	// user := user.User{UserName: "yo"}

// 	jsn, err := ioutil.ReadAll(r.Body)

// 	if err != nil {
// 		log.Fatal("Error reading the body", err)
// 	}

// 	err = json.Unmarshal(jsn, &user) // put info into the user struct
// 	user.ID = bson.NewObjectId()

	
// 	if err = repo.InsertUser(user); err != nil {
// 		log.Fatal("Error adding new register in database", err)
// 	} 
	

// 	// fmt.Fprintln(w, user);
	
// }

// func UpdateUserEndPoint(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "not implemented yet !")
// }

// func DeleteUserEndPoint(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "not implemented yet !")
// }


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", userapi.FindAll).Methods("GET")
	r.HandleFunc("/user/{id}", userapi.Find).Methods("GET")
	r.HandleFunc("/users", userapi.Create).Methods("POST")
	r.HandleFunc("/user/delete/{id}", userapi.Delete).Methods("DELETE")
	r.HandleFunc("/users", userapi.Update).Methods("PUT")

	corsObj := handlers.AllowedOrigins([]string{"*"})
	// r.HandleFunc("/users", CreateUserEndPoint).Methods("POST")
	// r.HandleFunc("/users", UpdateUserEndPoint).Methods("PUT")
	// r.HandleFunc("/users", DeleteUserEndPoint).Methods("DELETE")
	if err := http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}), corsObj)(r)); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Listening on port 3000....")
	}
}