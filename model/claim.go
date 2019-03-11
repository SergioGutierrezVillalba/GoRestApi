package model

import (
	
	jwt "github.com/dgrijalva/jwt-go"

	"FirstProject/entities"
)

type Claim struct {
	entities.User `json:"user"`
	jwt.StandardClaims

}