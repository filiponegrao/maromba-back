package router

import (
	"log"

	"github.com/filiponegrao/maromba-back/controllers"

	"github.com/gin-gonic/gin"
)

// Initialize : Initialize
func Initialize(r *gin.Engine) {
	// the jwt middleware
	authMiddleware, err := controllers.GetAuthMiddlware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	r.Use(gin.Recovery())
	r.Use(controllers.CORSMiddleware())
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(400, nil)
	// })
	api := r.Group("/api")

	// Sem autenticacao
	api.POST("/users", Logger(), controllers.CreateUser)
	api.POST("/login", Logger(), authMiddleware.LoginHandler)
	api.POST("/refresh", Logger(), controllers.CheckRefreshCode, authMiddleware.RefreshHandler)
	api.POST("/users/password/forgot", Logger(), controllers.ForgotPassword)

	// Com autenticacao
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.Use(Logger())
		api.GET("/logged", controllers.RequestUserLogged)
		api.Use(Authorizer())
		api.POST("/users/password/change", controllers.ChangePassword)
		api.PUT("/users", controllers.UpdateUser)
		api.GET("/logs", controllers.GetLogs)
		api.GET("/dashboard/users/list", controllers.GetUsers)
	}
}
