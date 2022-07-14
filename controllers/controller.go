package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/filiponegrao/maromba-back/config"
	dbpkg "github.com/filiponegrao/maromba-back/db"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var Conf config.Configuration

func MainConfig(config config.Configuration) {
	Conf = config
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Refresh-Code, Application-Version")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
	}
}

// GetUserLogged is responsbile both for return the
// logged user, and send an error to the requester
// if is not logged
func GetUserLogged(c *gin.Context) models.User {
	db := dbpkg.DBInstance(c)
	var user models.User
	claims := jwt.ExtractClaims(c)
	userIDString := claims["id"]
	if userIDString == nil {
		return models.User{}
	}
	userID := int64(userIDString.(float64))
	if err := db.First(&user, userID).Error; err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{"error": "Usuário inválido"})
		c.Abort()
	}
	return user
}

// RespondError is responsible to forward all error messages
// to the requester.
func RespondError(c *gin.Context, errString string, code int) {
	SaveLogData(c, code, errString)
	if IsVersion2(c) {
		RespondErrorV2(c, errString, 400, code)
	} else {
		c.JSON(code, gin.H{"error": errString})
	}
}

func RespondErrorV2(c *gin.Context, errString string, httpstatus int, code int) {
	resoonse := models.RequestResponse{
		Success: false,
		Message: errString,
		Result:  code,
	}
	c.JSON(code, resoonse)
}

// RespondSuccess is responsible to forward all success messages
// to the requester.
func RespondSuccess(c *gin.Context, content interface{}) {
	data, err := json.Marshal(content)
	if err != nil {
		log.Println(err)
	}
	dataString := ""
	if len(data) > 0 {
		dataString = string(data)
	}
	SaveLogData(c, 200, dataString)
	c.JSON(200, content)
}

func RespondSuccessV2(c *gin.Context, content interface{}) {
	SaveLogData(c, 200, "")
	resoonse := models.RequestResponse{
		Success: true,
		Message: "",
		Result:  content,
	}
	c.JSON(200, resoonse)
}

func SaveLogData(c *gin.Context, responseCode int, message string) {
	db := dbpkg.DBInstance(c)
	object, _ := c.Get("log")
	if object == nil {
		return
	}
	logData := object.(models.Log)
	logData.ResponseCode = responseCode
	logData.ResponseBody = message
	db.Save(&logData)
}

func GetAmericanDateStringFrom(dateString string) string {
	parts := strings.Split(dateString, "/")
	newDateString := parts[2] + "-" + parts[1] + "-" + parts[0]
	return newDateString
}

func UserSearchTermQueryDb(c *gin.Context, db *gorm.DB) (result *gorm.DB) {
	searchTerm := c.Query("searchTerm")
	if searchTerm != "" {
		searchTerms := strings.Split(searchTerm, " ")
		for index, term := range searchTerms {
			searchTermQuery := "%" + term + "%"
			if index == 0 {
				db = db.Where(
					"name like ? or email like ? or cpf like ? or cnpj like ?",
					searchTermQuery, searchTermQuery, searchTermQuery, searchTermQuery,
				)
			} else {
				db = db.Or(
					"name like ? or email like ? or cpf like ? or cnpj like ?",
					searchTermQuery, searchTermQuery, searchTermQuery, searchTermQuery,
				)
			}
		}
	}
	return db
}

func IsVersion2(c *gin.Context) bool {
	if c.Request.Header.Get("Application-Version") == "v1" {
		return true
	} else {
		return false
	}
}

// MARK: Auxiliary

func setMonthDateQueryString(monthDateString string, db *gorm.DB) (dbResult *gorm.DB, err error) {
	monthDateStringParts := strings.Split(monthDateString, "/")
	if len(monthDateStringParts) != 2 {
		return db, errors.New("Data incorreta")
	}
	if len(monthDateStringParts[0]) > 2 {
		return db, errors.New("Data incorreta")
	}
	if len(monthDateStringParts[1]) != 4 {
		return db, errors.New("Data incorreta")
	}
	month, err := strconv.Atoi(monthDateStringParts[0])
	if err != nil {
		return db, err
	}
	year, err := strconv.Atoi(monthDateStringParts[1])
	if err != nil {
		return db, err
	}
	now := time.Now()
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
	finalDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, now.Location())

	dbResult = db.Where("date > ? and date < ?", startDate, finalDate)
	return dbResult, nil
}

func CreateExternalLog(
	path string,
	body string,
	response string,
	headers string,
	method string,
	code int,
	user models.User,
	c *gin.Context) (log models.Log, err error) {

	log.UserID = user.ID
	log.UserMail = user.Email
	log.UserName = user.Name
	log.Method = c.Request.Method
	log.ResponseBody = response
	log.ResponseCode = code
	log.Body = body
	log.URL = c.Request.Host + c.Request.URL.Path
	log.Header = headers
	now := time.Now()
	log.CreatedAt = &now
	log.Description = log.GetDescriptionString()
	log.IP = c.ClientIP()
	log, err = SaveLog(c, log)
	return log, err
}
