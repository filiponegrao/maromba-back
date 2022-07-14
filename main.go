package main

import (
	"io"
	"log"
	"os"

	"github.com/filiponegrao/maromba-back/config"
	"github.com/filiponegrao/maromba-back/controllers"
	"github.com/filiponegrao/maromba-back/db"
	"github.com/filiponegrao/maromba-back/server"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var conf config.Configuration

func init() {}

func GetConfig(path string) {
	conf = config.Get(path)
	controllers.ConfigEmail(conf)
	controllers.MainConfig(conf)
	db.SetConfigurations(conf)
	log.Println("Arquivo de configuração " + path + " lido com sucesso!")
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		log.Fatal("Faltando arquivo de configuração json")
		return
	}
	configFileName := os.Args[1]
	GetConfig(configFileName)
	f, err := os.OpenFile(conf.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	configDatabaseLog(database, f)

	s := server.Setup(database)
	port := conf.ApiPort
	s.Run(":" + port)
}

func configDatabaseLog(database *gorm.DB, f *os.File) {
	logger := gorm.Logger{
		LogWriter: log.New(io.MultiWriter(f, os.Stdout), "\r\n", log.LstdFlags),
	}
	database.SetLogger(logger)
}
