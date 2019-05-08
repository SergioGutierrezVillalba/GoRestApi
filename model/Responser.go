package model

import (
	"encoding/json"
	"log"
	"net/http"
)

type Responser struct{}

func (responser *Responser) WithError(w http.ResponseWriter, code int, payload interface{}) {
	responser.WithJson(w, code, payload)
}

func (responser *Responser) WithJson(w http.ResponseWriter, code int, payload interface{}) {

	response, err := json.Marshal(payload)

	if err != nil {
		log.Fatal("Error al convertir la entidad a JSON")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}
