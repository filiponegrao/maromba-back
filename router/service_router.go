package router

import (
	"github.com/filiponegrao/maromba-back/controllers"
	"github.com/gin-gonic/gin"
)

func ServiceRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !controllers.Conf.Features.Service {
			controllers.RespondError(c, "Funcionalidade não habilitada.", 400)
			c.Abort()
		}
		c.Next()
	}
}
