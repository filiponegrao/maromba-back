package router

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/filiponegrao/maromba-back/controllers"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/filiponegrao/maromba-back/tools"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verifica se é requisicao de log
		if strings.Contains(c.Request.URL.String(), "api/logs") {
			c.Next()
			return
		}

		user := controllers.GetUserLogged(c)
		out := gin.DefaultWriter
		userLog := "\n********* New Request *********"
		userLog += "\n[Path: " + c.Request.Method + " " + c.Request.URL.String() + "]\n"
		userLog += "\n[User: " + strconv.FormatInt(user.ID, 10) + " " + user.Email + "]\n"

		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
		bodyString := readBody(rdr1)
		logFileBodyString := bodyString
		log := models.Log{}
		// Verifica se possui arquivo
		if strings.Contains(c.Request.URL.String(), "/upload") {
			logFileBodyString = "file" //strings.Split(logFileBodyString, "Content-type")[0]
		} else if strings.Contains(c.Request.URL.String(), "api/login") {
			// logFileBodyString = "credentials"
		}
		userLog += "\n[Body: " + logFileBodyString + "]\n"
		fmt.Fprint(out, userLog)

		log.UserID = user.ID
		log.UserMail = user.Email
		log.UserName = user.Name
		log.Method = c.Request.Method
		log.Body = logFileBodyString
		log.URL = c.Request.Host + c.Request.URL.Path
		log.Header = tools.HeadersToString(c.Request.Header)
		now := time.Now()
		log.CreatedAt = &now
		log.Description = log.GetDescriptionString()
		log.IP = c.ClientIP()
		logSaved, err := controllers.SaveLog(c, log)
		if err != nil {
			fmt.Fprint(out, "LOG ERROR: "+err.Error())
		}
		c.Request.Body = rdr2
		c.Set("log", logSaved)
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
