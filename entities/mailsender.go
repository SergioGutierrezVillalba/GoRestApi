package entities

import (
	"log"
	"net/smtp"
)

type MailSender struct{}

func (mailSender *MailSender)Send(email string, token string) error {

	from := "sendermail495@gmail.com"
	pass := "mail888mail"
	to := email
	urlToRecover := "http://localhost:4002/recover/" + token

	msg:= "From: "+ from + "\n" +
		  "To: "+ to + "\n" +
		  "Subject: Hola! \n\n"+
		  "Use this link to recover your password: " + urlToRecover
	
	err:= smtp.SendMail("smtp.gmail.com:587",
		   smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		   from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Print("sent")
	return err
	
}