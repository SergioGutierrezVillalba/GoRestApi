package jwt

import (
	"time"

	"FirstProject/Model"
	claims "FirstProject/Model/JWT/Claim"

	djwt "github.com/dgrijalva/jwt-go"
)

type JWTPayload struct {
	Claims		claims.Claim
}

var (
	StandardExpiration 	=  time.Now().Add(time.Hour * 1).Unix()
)
// TODO
func NewPayload(user model.User) JWTPayload {
	payload := JWTPayload{
		Claims: claims.Claim{
			User: user,
			StandardClaims: djwt.StandardClaims{
				ExpiresAt: StandardExpiration,
				Subject: user.Username,
			},
		},
	}
	return payload
}

func NewPayloadFromRaw() JWTPayload {
	payload := JWTPayload{
		Claims: claims.Claim{
			User: model.User{},
			StandardClaims: djwt.StandardClaims{
				ExpiresAt: 0,
				Subject: "Unknown",
			},
		},
	}
	return payload
}

func (j *JWTPayload) GetClaims() claims.Claim {
	return j.Claims
}

func (j *JWTPayload) SetClaims(myClaims *MyCustomClaims) {
	j.Claims = claims.Claim{
		User: myClaims.User,
		StandardClaims: myClaims.StandardClaims,
	}
}

func (j *JWTPayload) GetExpiration() int64 {
	return j.Claims.GetExpiration()
}

func (j *JWTPayload) GetSubject() string {
	return j.Claims.GetSubject()
}

// TODO Getters...