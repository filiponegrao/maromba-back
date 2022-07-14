package email

import (
	"github.com/filiponegrao/maromba-back/config"
	"github.com/filiponegrao/maromba-back/models"
)

func UserPendingInvite(conf config.Configuration, invite models.Invite, code string) error {
	subject := "[" + conf.ApiName + "]: Conta criada!"
	message := "Olá " + invite.Invited.Name + "!"
	message += "\n\nSua conta foi criada com sucesso!"
	message += "\n\nPara ativar sua conta insira o código de acesso quando requisitado."
	message += "\n* Código de Acesso: " + code
	message += "\n\nEquipe " + conf.ApiName
	message += "\n\nPara mais informações: " + conf.Site
	return emailConf.sendEmail(conf, invite.Invited.Email, message, subject, 0)

}
