package model

import (
	"net/http"
	"encoding/json"
	"log"
)

type Responser struct {}

func (responser *Responser) WithError(w http.ResponseWriter, code int, msg string){
	responser.WithJson(w, code, map[string]string{"error": msg})
}

func (responser *Responser) WithJson(w http.ResponseWriter, code int, payload interface{}){

	response, err := json.Marshal(payload)

	if err != nil {
		log.Fatal("Error al convertir la entidad a JSON")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}