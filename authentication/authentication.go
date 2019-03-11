package authentication

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	"FirstProject/model"
	"FirstProject/entities"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	privateBytes, err := ioutil.ReadFile("./private.rsa")
	if err != nil {
		log.Fatal("No se ha podido leer la clave privada")
	}

	publicBytes, err := ioutil.ReadFile("./public.rsa.pub")
	if err != nil {
		log.Fatal("No se ha podido leer la clave pública")
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se ha podido parsear la clave privada")
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("No se ha podido parsear la clave pública")
	}
}

func GenerateJWT(user entities.User) string {
	claims := model.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer: "App usuarios",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	
	if err != nil {
		log.Fatal("No se ha podido firmar el token")
	} 

	return result
}