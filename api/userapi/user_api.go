package userapi

import (
	"FirstProject/config"
	"FirstProject/model"
	"FirstProject/entities"

	"fmt"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2/bson"
)



func FindAll(response http.ResponseWriter, request *http.Request){
	db, err := config.Connect()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return 

	} else {
		userModel := model.UserModel{
			Db: db,
		}

		users, err2 := userModel.FindAll()

		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err2.Error())
			return 
		} else {
			respondWithJson(response, http.StatusOK, users)
		}
	}
}

func Find(response http.ResponseWriter, request *http.Request){
	db, err := config.Connect()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return 

	} else {
		userModel := model.UserModel{
			Db: db,
		}

		vars := mux.Vars(request)
		id := vars["id"]

		user, err2 := userModel.Find(id)

		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err2.Error())
			return 
		} else {
			respondWithJson(response, http.StatusOK, user)
		}
	}
}

func Create(response http.ResponseWriter, request *http.Request){

	db, err := config.Connect()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return 

	} else {
		userModel := model.UserModel{
			Db: db,
		}
		
		var user entities.User
		user.Id = bson.NewObjectId() // generates new id in Bson notation
		json.NewDecoder(request.Body).Decode(&user) // transform user struct into JSON notation

		err2 := userModel.Create(&user)

		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err2.Error())
			return 
		} else {
			respondWithJson(response, http.StatusOK, user)
		}
	}
}

func Delete(response http.ResponseWriter, request *http.Request){

	db, err := config.Connect()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return 

	} else {
		userModel := model.UserModel{
			Db: db,
		}
		
		vars := mux.Vars(request)
		id := vars["id"]

		err2 := userModel.Delete(id)

		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err2.Error())
			return 
		} else {
			respondWithJson(response, http.StatusOK, nil)
			fmt.Println("Deleted user with id: " + id)
		}
	}
}

func Update(response http.ResponseWriter, request *http.Request){

	db, err := config.Connect()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return 

	} else {
		userModel := model.UserModel{
			Db: db,
		}

		var user entities.User
		json.NewDecoder(request.Body).Decode(&user)

		err2 := userModel.Update(&user)

		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err2.Error())
			return 
		} else {
			respondWithJson(response, http.StatusOK, user)
		}
	}
}



func respondWithError(w http.ResponseWriter, code int, msg string){
	respondWithJson(w, code, map[string]string{"error":msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}){
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}