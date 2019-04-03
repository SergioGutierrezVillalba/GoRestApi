package auth

import (
	"crypto/rsa"

	"io/ioutil"
	"log"
	"fmt"
	"time"

	"FirstProject/model"

	djwt "github.com/dgrijalva/jwt-go"
	alyu "github.com/alyu/encrypt"
)

type Authentication struct{}
type MyCustomClaims struct{
	djwt.StandardClaims
	User 	model.User `json:"user"`
}

var (
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
	passphrase = "1234567891011121"
	secret 	   = "gnommoSecretKey@#!jejejej6567606" // 32 bytes
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

	privateKey, err = djwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se ha podido parsear la clave privada")
	}

	publicKey, err = djwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("No se ha podido parsear la clave pública")
	}
}

func (a *Authentication) GenerateJWT(user model.User) string {

	user.RawId = user.Id.Hex()

	claims := Claim{
		User: user,
		StandardClaims: djwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer: "App Timer",
		},
	}

	token := djwt.NewWithClaims(djwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)

	if err != nil {
		log.Fatal("No se ha podido firmar el token")
	} 
	return a.EncodeBase64(a.Encrypt(result))
}

func (a *Authentication) GetClaimsFromJwt(rawJWT string) (djwt.Claims, error) {
	token, err := djwt.ParseWithClaims(rawJWT, &MyCustomClaims{}, func(token *djwt.Token)(interface{}, error){
		return []byte(privateKey.D.Bytes()), nil
	})
	return token.Claims, err
}

func (a *Authentication) GetUserIdFromJWT(rawJWT string) (string, error) {
	jwt := a.Decrypt(a.DecodeBase64(rawJWT))
	fmt.Println(rawJWT)
	claims, err := a.GetClaimsFromJwt(jwt)
	claimsCustomized := claims.(*MyCustomClaims)

	fmt.Println("ey: " + claimsCustomized.User.RawId)
	return claimsCustomized.User.RawId, err
}

func (a *Authentication) GetExpirationTimeOfJWT(rawJWT string) (int64){
	claims, _ := a.GetClaimsFromJwt(rawJWT)
	claimsCustomized := claims.(*MyCustomClaims)
	return claimsCustomized.StandardClaims.ExpiresAt
}

func (a *Authentication) IsNotValid(expiration int64) bool {
	if expiration > time.Now().Unix(){
		return false
	}
	return true
}

// func (a *Authentication) Hash(data string) (hashedData string, err error) {
// 	h := hmac.New(sha256.New, []byte(secret))
// 	h.Write([]byte(data))
// 	hashedData = hex.EncodeToString(h.Sum(nil))
// 	return
// }

func (a *Authentication) Encrypt(rawData string) (encryptedData []byte) {
	encryptedData = alyu.Encrypt([]byte(secret), []byte(rawData))
	return
}

func (a *Authentication) Decrypt(encryptedData []byte) (decryptedData string) {
	decryptedData = string(alyu.Decrypt([]byte(secret), encryptedData))
	return
}

func (a *Authentication) EncodeBase64(data []byte) (encodedData string){
	return alyu.EncodeBase64(data)
} 

func (a *Authentication) DecodeBase64(data string) (decodedData []byte){
	return alyu.DecodeBase64(data)
}


