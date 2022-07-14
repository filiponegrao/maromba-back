package controllers

import (
	"log"
	"net/mail"
	"net/smtp"

	"github.com/filiponegrao/maromba-back/config"
	"github.com/scorredoira/email"
)

var emailConf EmailConfiguration

type EmailConfiguration struct {
	Mail1     string
	Password1 string
	Server1   string
	Port1     string
	Mail2     string
	Password2 string
	Server2   string
	Port2     string
	Site      string
}

func ConfigEmail(conf config.Configuration) {
	confEmail := EmailConfiguration{
		Mail1:     conf.Email1,
		Password1: conf.EmailPassword1,
		Server1:   conf.EmailServer1,
		Port1:     conf.EmailPort1,
		Mail2:     conf.Email2,
		Password2: conf.EmailPassword2,
		Server2:   conf.EmailServer2,
		Port2:     conf.EmailPort2,
		Site:      conf.Site,
	}
	ConfigEmailEngine(confEmail)
}

func ConfigEmailEngine(engine EmailConfiguration) {
	emailConf = engine
}

func (emailConf EmailConfiguration) sendEmail(targetEmail string, text string, subject string, option int) error {

	emailString := ""
	password := ""
	server := ""
	port := ""
	if option <= 0 {
		emailString = emailConf.Mail1
		password = emailConf.Password1
		server = emailConf.Server1
		port = emailConf.Port1
	} else {
		emailString = emailConf.Mail2
		password = emailConf.Password2
		server = emailConf.Server2
		port = emailConf.Port2
	}
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		emailString,
		password,
		server,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	address := server + ":" + port

	message := email.NewMessage(subject, text)
	sender := mail.Address{}
	sender.Address = emailString
	sender.Name = "Suporte " + Conf.ApiName
	message.From = sender
	message.To = []string{targetEmail}

	// err := message.Attach("logo1.png")
	// if err != nil {
	// 	log.Println(err)
	// }

	err := email.Send(address, auth, message)
	if err != nil {
		log.Println("Erro ao tentar enviar email: " + err.Error())
	}

	return err
}
