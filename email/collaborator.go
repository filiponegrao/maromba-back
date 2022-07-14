package email

import (
	"github.com/filiponegrao/maromba-back/config"
	"github.com/filiponegrao/maromba-back/models"
)

func CollaboratorInvite(conf config.Configuration, invite models.Invite, password string, code string) error {
	subject := "[" + conf.ApiName + "]: Convite para a plataforma!"
	message := "Olá " + invite.Invited.Name + "!"
	message += "\n\nVoce foi convidado por " + invite.Inviter.Name + " para participar da plataforma "
	message += conf.ApiName + " como " + invite.Invited.GetUserTypeString() + "."
	message += "\n\nSua conta foi criada com as seguintes credenciais"
	message += "\n* Email: " + invite.Invited.Email + "\n* Senha: " + password
	message += "\n\nOBS: Recomendamos fortmenente a alteração desta senha logo após o primeiro login."
	message += "\n\nPara ativar sua conta faça login no site e insira o código de acesso quando requisitado."
	message += "\n* Site: " + conf.AdminURL
	message += "\n* Código de Acesso: " + code
	message += "\n\nEquipe " + conf.ApiName
	message += "\n\nPara mais informações: " + conf.Site
	return emailConf.sendEmail(conf, invite.Invited.Email, message, subject, 0)
}

func ReSendInvite(conf config.Configuration, invite models.Invite, code string) error {
	subject := "[" + conf.ApiName + "]: Convite para a plataforma!"
	message := "Olá " + invite.Invited.Name + "!"
	message += "\n\nVoce foi convidado por " + invite.Inviter.Name + " para participar da plataforma "
	message += conf.ApiName + " como " + invite.Invited.GetUserTypeString() + "."
	message += "\n\nPara ativar sua conta faça login no site e insira o código de acesso quando requisitado."
	message += "\n* Site: " + conf.AdminURL
	message += "\n* Código de Acesso: " + code
	message += "\n\nEquipe " + conf.ApiName
	message += "\n\nPara mais informações: " + conf.Site

	return emailConf.sendEmail(conf, invite.Invited.Email, message, subject, 0)
}
