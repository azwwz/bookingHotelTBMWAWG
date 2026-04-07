package main

import (
	"github.com/azwwz/bookingHotelTBMWAWG/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"os"
	"strings"
	"time"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {

	// set server
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// get a client from sever connect
	client, err := server.Connect()
	if err != nil {
		log.Println(err)
	}

	// generate a mail
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)

	} else {
		file, err := os.ReadFile("./email-templates/" + m.Template)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}
		fileString := string(file)
		fileString = strings.Replace(fileString, "[%body%]", m.Content, -1)
		email.SetBody(mail.TextHTML, fileString)
	}

	// combine a mail and client
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("email sent")
	}

}
