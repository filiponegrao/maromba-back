package email

import "github.com/filiponegrao/maromba-back/config"

func EmailChangedPassword(conf config.Configuration, targetEmail string) error {

	message := conf.ApiName + " informa:\n\n"
	message += "Sua senha foi alterada com sucesso!\n\n"
	message += "Obs: Se voce nao solicitou uma troca de senha "
	message += "entre em contato no site: \n"
	message += emailConf.Site
	message += "\n\nAtenciosamente,\n"
	message += "Equipe Acelerados"

	subject := "[Acelerados]: Senha alterada com sucesso!"

	return emailConf.sendEmail(conf, targetEmail, message, subject, 0)
}

func EmailPasswordNew(conf config.Configuration, targetEmail string, password string) error {

	message := conf.ApiName + " informa:\n\n"
	message += "Sua nova senha foi gerada!\n"
	message += "senha: " + password + "\n\n"
	message += "Obs: Se voce nao solicitou uma troca de senha "
	message += "entre em contato no site: \n"
	message += emailConf.Site
	message += "\n\nAtenciosamente,\n"
	message += "Equipe " + conf.ApiName

	subject := "[" + conf.ApiName + "]: Nova senha!"

	return emailConf.sendEmail(conf, targetEmail, message, subject, 0)
}
