package jwt

import (
	"log"
	"time"
	"fmt"
	"io/ioutil"
	"crypto/rsa"

	"FirstProject/Model"

	djwt "github.com/dgrijalva/jwt-go"
	alyu "github.com/alyu/encrypt"
)

type JWT struct {
	Header		JWTHeader
	Payload		JWTPayload
	Expiration	int64
	Raw			string
	Encrypted	[]byte
	Encoded		string
	Decoded		[]byte
	Decrypted	string
}

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

// NewJWT
// SignTokenRS256()
// Encrypt()
// rawJWT = EncodeBase64()

// decodedJWT = DecodeBase64()
// DecryptDecoded()
// DumpDataInsideJWTStruct()

func NewJWT(user model.User) JWT {
	jwt := JWT{
		Header: NewHeader(),
		Payload: NewPayload(user),
	}
	return jwt
}

func NewJWTFromRaw(rawJWT string) JWT {
	jwt := JWT{
		Header: NewHeaderFromRaw(),
		Payload: NewPayloadFromRaw(),
		Raw: rawJWT,
	}
	return jwt
}

func (j *JWT) SignTokenRS256(){
	token := djwt.NewWithClaims(djwt.SigningMethodRS256, &j.Payload.Claims)
	result, err := token.SignedString(privateKey)

	if err != nil {
		log.Fatal("No se ha podido firmar el token")
	}

	j.Raw = result
}

func (j *JWT) GetExpiration() int64 {
	return j.Payload.GetExpiration()
}

func (j *JWT) IsExpirated() bool {
	if j.GetExpiration() > time.Now().Unix(){
		return false
	}
	return true
}

func (j *JWT) Encrypt(){
	j.Encrypted = alyu.Encrypt([]byte(secret), []byte(j.Raw))
}

// Obligation of encrypt the token before encoding.
func (j *JWT) EncodeBase64() string {
	j.Encoded = alyu.EncodeBase64(j.Encrypted)
	return j.Encoded
}

// Obligation of decoding the token before decrypting.
func (j *JWT) DecodeBase64() (){
	j.Decoded = alyu.DecodeBase64(j.Raw)
}

func (j *JWT) Decrypt(){
	j.Decrypted = string(alyu.Decrypt([]byte(secret), j.Encrypted))
}

func (j *JWT) DecryptDecoded(){
	j.Decrypted = string(alyu.Decrypt([]byte(secret), j.Decoded))
	fmt.Println("-----------------------")
	fmt.Println(j.Decrypted)
	fmt.Println("-----------------------")
}

func (j *JWT) DumpDataInsideJWTPayload(){

	claims, err := j.GetClaims()
	myClaims := claims.(*MyCustomClaims)

	if err != nil {
		log.Panic(err)
	}

	j.Payload.SetClaims(myClaims)
}

func (j *JWT) GetClaims() (djwt.Claims, error) {
	token, err := djwt.ParseWithClaims(j.Decrypted, &MyCustomClaims{}, func(token *djwt.Token)(interface{}, error){
		return []byte(privateKey.D.Bytes()), nil
	})
	return token.Claims, err
}