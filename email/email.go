package email

import (
	"errors"
	"log"
	"net/mail"
	"net/smtp"
	"strings"

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

func ConfigEmailEngine(engine EmailConfiguration) {
	emailConf = engine
}

func (emailConf EmailConfiguration) sendEmail(conf config.Configuration, targetEmail string, text string, subject string, option int) error {

	emailString := ""
	password := ""
	server := ""
	port := ""
	if option <= 0 {
		emailString = conf.Email1
		password = conf.EmailPassword1
		server = conf.EmailServer1
		port = conf.EmailPort1
	} else {
		emailString = conf.Email2
		password = conf.EmailPassword2
		server = conf.EmailServer2
		port = conf.EmailPort2
	}
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		emailString,
		password,
		server,
	)
	address := server + ":" + port
	println("target")
	println(address)
	message := email.NewMessage(subject, text)
	sender := mail.Address{}
	sender.Address = emailString
	sender.Name = "Suporte " + conf.ApiName
	message.From = sender
	message.To = []string{targetEmail}
	err := email.Send(address, auth, message)
	if err != nil {
		if strings.Contains(err.Error(), "503 need RCPT before DATA") {
			return errors.New("E-mail inválido!")
		}
		log.Println("Erro ao tentar enviar email: " + err.Error())
	}

	return err
}
