package entities

import ()

type PasswordRecover struct {
	NewPassword	   string 	`bson:"password" json:"password"`
	Token		   string 	`bson:"token" json:"token"`
}