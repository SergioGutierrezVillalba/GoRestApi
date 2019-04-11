package mail

import (
	"log"
	"bytes"
	"strings"
	"os"
	"fmt"
	"net/smtp"
	"html/template"
	
	times "FirstProject/Model/times"
	model "FirstProject/Model"
	gomail "gopkg.in/gomail.v2"
)

type MailSender struct {}
type ResponseEmail struct {
	Start			string
	Finish			string
	StartDate		string
	FinishDate		string
}

var (
	timePresenter times.TimePresenter
)

func (m *MailSender)SendRecover(email string, token string) error {

	from := "sendermail495@gmail.com"
	pass := "mail888mail"
	to := email
	urlToRecover := "http://localhost:4200/recover/" + token

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

func (ms *MailSender)SendFinishedTime(email string, timer model.TimerFormatted) error {

	projectPath := os.Getenv("GO_PROJECT_CONF_ROUTE")
	relativePath := "Templates/"
	templatesPath := projectPath + relativePath
	fmt.Println("Ruta: " + templatesPath)

	startInfoSliced := strings.Fields(timer.Start)
	finishInfoSliced := strings.Fields(timer.Finish)

	// INITTIME
	startDate := startInfoSliced[0]
	startTime := startInfoSliced[1]

	// FINISHTIME
	finishDate := finishInfoSliced[0]
	finishTime := finishInfoSliced[1]

	// DURATION-------------------------
	fmt.Println("Duration(s):")
	fmt.Println(timer.Duration)

	duration := timePresenter.SecondsToHuman(int(timer.Duration))

	fmt.Println("StartDate : " + startDate + "|" + "StartTime: " + startTime)
	fmt.Println("FinishDate : " + finishDate + "|" + "FinishTime: " + finishTime)

	from := "sendermail495@gmail.com"
	pass := "mail888mail"
	to := email

	t := template.New("email-template.html")

	var err error
	t, err = t.ParseFiles(templatesPath + "email-template.html")
	if err != nil {
		fmt.Println(err)
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, struct{
			Start			string
			Finish			string
			StartDate		string
			FinishDate		string
			Duration		string
		}{
			Start:			startTime,
			Finish:			finishTime,
			StartDate:		startDate,
			FinishDate:		finishDate,	
			Duration:		duration,	
		})

	if err != nil {
		fmt.Println(err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Timer finished")
	m.SetBody("text/html", result)
	// m.Attach("template.html") attach whatever you want

	d := gomail.NewDialer("smtp.gmail.com", 587, from, pass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	log.Print("sent")
	return err	
}