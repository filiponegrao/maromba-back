package server

import (
	"github.com/filiponegrao/maromba-back/middleware"
	"github.com/filiponegrao/maromba-back/router"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.SetDBtoContext(db))
	router.Initialize(r)
	return r
}
