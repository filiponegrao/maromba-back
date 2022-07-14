package controllers

import (
	"log"
	"time"

	dbpkg "github.com/filiponegrao/maromba-back/db"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/filiponegrao/maromba-back/tools"
	"github.com/gin-gonic/gin"
)

func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	var loginResponse loginResponse
	loginResponse.Token = token
	loginResponse.Expire = &expire
	// Gera um refresh code
	refreshCode, err := CreateRefreshCode(c, token)
	if err != nil {
		log.Println(err.Error())
		RespondError(c, "Erro interno!", 400)
		c.Abort()
		return
	}
	loginResponse.RefreshCode = refreshCode.Code
	RespondSuccess(c, loginResponse)
}

func CheckRefreshCode(c *gin.Context) {
	db := dbpkg.DBInstance(c)

	refreshCodes := c.Request.Header["Refresh-Code"]
	if len(refreshCodes) == 0 {
		message := "Sem acesso 1"
		RespondError(c, message, 401)
		c.Abort()
		return
	}
	refreshCode := refreshCodes[0]
	authorizations := c.Request.Header["Authorization"]
	if len(authorizations) == 0 {
		message := "Sem acesso 2"
		RespondError(c, message, 401)
		c.Abort()
		return
	}
	authorization := authorizations[0]
	log.Println(authorization)
	var refreshCodeSaved models.RefreshCode
	if err := db.
		Where("code = ? AND access_token = ? AND status = 0", refreshCode, authorization).
		First(&refreshCodeSaved).Error; err != nil {
		message := "Sem acesso 3"
		RespondError(c, message, 401)
		c.Abort()
		return
	}
	// Inutiliza esse código
	refreshCodeSaved.Status = 1
	if err := db.Save(&refreshCodeSaved).Error; err != nil {
		message := "Erro interno"
		RespondError(c, message, 400)
		c.Abort()
		return
	}

}

func CreateRefreshCode(c *gin.Context, accessToken string) (models.RefreshCode, error) {
	db := dbpkg.DBInstance(c)
	now := time.Now()
	refreshCodeString := tools.RandomString(6)
	refreshCodeString = tools.EncryptTextSHA512(refreshCodeString)
	var refreshCode models.RefreshCode
	refreshCode.Code = refreshCodeString
	refreshCode.CreatedAt = &now
	refreshCode.AccessToken = "Bearer " + accessToken
	if err := db.Save(&refreshCode).Error; err != nil {
		log.Println(err.Error())
		return refreshCode, err
	}
	return refreshCode, nil
}
