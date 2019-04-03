package auth

import (
	
	jwt "github.com/dgrijalva/jwt-go"

	"FirstProject/model"
)

type Claim struct {
	model.User `json:"user"`
	jwt.StandardClaims

}