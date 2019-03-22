package controllers

import (
	"FirstProject/model"
	auth "FirstProject/model/auth"
)

var (
	respond 		model.Responser
	mailSender  	model.MailSender

	crypter     	auth.Crypter
	requestHeaders 	auth.RequestHeaders
	authenticator	auth.Authentication
	responseToken	auth.ResponseToken

	// checker 		validation.Checker
)