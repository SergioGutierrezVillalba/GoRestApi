package userapi

import (
	"FirstProject/config"
	"FirstProject/entities"
	"FirstProject/model"
	auth "FirstProject/authentication"

	"fmt"
	"log"
	"strings"
	"net/http"
	// "io/ioutil"

	"encoding/json"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type UserAPI struct {
}

func (userApi *UserAPI) FindAll(response http.ResponseWriter, request *http.Request) {

	db, err := config.Connect()
	defer db.Session.Close()
	
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

func (userApi *UserAPI) Find(response http.ResponseWriter, request *http.Request) {

	db, err := config.Connect()
	defer db.Session.Close()

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

func (userApi *UserAPI) Create(response http.ResponseWriter, request *http.Request) {

	db, err := config.Connect()
	defer db.Session.Close()

	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return

	} else {
		userModel := model.UserModel{
			Db: db,
		}

		var user entities.User
		user.Id = bson.NewObjectId()                // generates new id in Bson notation
		json.NewDecoder(request.Body).Decode(&user) // transform user struct into JSON notation

		if len(user.Username) > 0 && len(user.Password) > 0 {

			err2 := userModel.Create(&user)

			if err2 != nil {
				respondWithError(response, http.StatusBadRequest, err2.Error())
				return
			} else {
				respondWithJson(response, http.StatusOK, user)
			}

		} else {
			respondWithError(response, http.StatusBadRequest, "Campos necesarios están vacíos")
		} 
		
	}
}

func (userApi *UserAPI) Delete(response http.ResponseWriter, request *http.Request) {

	db, err := config.Connect()
	defer db.Session.Close()

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

func (userApi *UserAPI) Update(response http.ResponseWriter, request *http.Request) {

	db, err := config.Connect()
	defer db.Session.Close()

	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return

	} else {
		userModel := model.UserModel{
			Db: db,
		}

		var user entities.User
		json.NewDecoder(request.Body).Decode(&user)

		if len(user.Username) > 0 && len(user.Password) > 0 && len(user.Id) > 0 {
			
			err2 := userModel.Update(&user)

			if err2 != nil {
				respondWithError(response, http.StatusBadRequest, err2.Error())
				return
			} else {
				respondWithJson(response, http.StatusOK, user)
			}

		} else {
			respondWithError(response, http.StatusBadRequest, "Campos necesarios están vacíos")
		}

		
	}
}

func (userApi *UserAPI) Login(response http.ResponseWriter, request *http.Request){

	db, err := config.Connect()
	defer db.Session.Close()

	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
		return

	} else {
		userModel := model.UserModel{
			Db: db,
		}

		var user entities.User
		json.NewDecoder(request.Body).Decode(&user)

		generateJWT, err2 := userModel.Login(&user)

		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err2.Error())
			return
		} else {
			if generateJWT == "yes" {
				token := auth.GenerateJWT(user)
				result := model.ResponseToken{token}
				respondWithJson(response, http.StatusOK, result)
				
			} else {
				respondWithError(response, http.StatusBadRequest, "User does not exist")
			}
		}
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {

	response, err := json.Marshal(payload)

	if err != nil {
		log.Fatal("Error al convertir la entidad a JSON (user_api.go)")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

/*func fieldsAreEmpty(r http.Request) error {

	bodyReq, err := ioutil.ReadAll(r.Body)
	var user entities.User 

	if err != nil {
		return err
	} else {
		json.Unmarshal(bodyReq, &user)
		fmt.Printf("%+v\n", user.Username)
		return nil
	}
}*/

func formatRequest(r *http.Request) string {

	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
	  name = strings.ToLower(name)
	  for _, h := range headers {
		request = append(request, fmt.Sprintf("%v: %v", name, h))
	  }
	}
	
	// If this is a POST, add post data
	if r.Method == "POST" {
	   r.ParseForm()
	   request = append(request, "\n")
	   request = append(request, r.Form.Encode())
	} 
	 // Return the request as a string
	return strings.Join(request, "\n")
   }
