package router

import (
	"github.com/filiponegrao/maromba-back/controllers"
	"github.com/gin-gonic/gin"
)

func PetRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !controllers.Conf.Features.PetModule {
			controllers.RespondError(c, "Funcionalidade não habilitada.", 400)
			c.Abort()
		}
		c.Next()
	}
}
