package router

import (
	"github.com/filiponegrao/maromba-back/controllers"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/gin-gonic/gin"
)

func Authorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := controllers.GetUserLogged(c)
		// Verifica se está pendente
		if user.Status == models.USER_STATUS_PENDING {
			message := "Necessário confirmar a conta!"
			controllers.RespondError(c, message, 403)
			c.Abort()
		}
		// Verifica se está desativado
		if user.Status == models.USER_STATUS_BLOCKED {
			message := "Sem acesso ao aplicativo! Consulte o suporte!"
			controllers.RespondError(c, message, 403)
			c.Abort()
		}
		c.Next()
	}
}
