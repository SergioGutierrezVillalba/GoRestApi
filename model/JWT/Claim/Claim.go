package claim

import (
	
	jwt "github.com/dgrijalva/jwt-go"

	"FirstProject/Model"
)

type Claim struct {
	User 			model.User `json:"user"`
	StandardClaims 	jwt.StandardClaims
}

func (c *Claim) GetExpiration() int64 {
	return c.StandardClaims.ExpiresAt
} 

func (c *Claim) GetSubject() string {
	return c.StandardClaims.Subject
}

func (c *Claim) GetAudience() string {
	return c.StandardClaims.Audience
}

func (c *Claim) GetId() string {
	return c.StandardClaims.Id
}

func (c *Claim) GetIssuer() string {
	return c.StandardClaims.Issuer
}

func (c *Claim) GetIssuedAt() int64 {
	return c.StandardClaims.IssuedAt
}

func (c *Claim) Valid() error {
	return nil
}
