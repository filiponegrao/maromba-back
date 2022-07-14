package db

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/filiponegrao/maromba-back/config"
	"github.com/filiponegrao/maromba-back/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/serenize/snaker"
)

var conf config.Configuration

func SetConfigurations(configuration config.Configuration) {
	conf = configuration
}

func Connect() (*gorm.DB, error) {
	// Set defaul database
	database := "sqlite3"
	// Get the input database
	if conf.Database != "" {
		database = conf.Database
	}
	var db *gorm.DB
	var err error
	if database == "postgres" {
		log.Println("Utilizando conexão com o postgresql...")
		path := "host=" + conf.DatabaseIp + " port=" + conf.DatabasePort
		path += " user=" + conf.DatabaseUsername + " dbname=" + conf.DatabaseName
		path += " password=" + conf.DatabasePassword
		db, err = gorm.Open("postgres", path)
	} else if database == "sqlite3" {
		log.Println("Utilizando conexão com o sqlite3...")
		dir := filepath.Dir("db/database.db")
		db, err = gorm.Open("sqlite3", dir+"/database.db")
	} else {
		log.Println("Utilizando conexão default com o sqlite3...")
		// Uses sqlite3 as default
		dir := filepath.Dir("db/database.db")
		db, err = gorm.Open("sqlite3", dir+"/database.db")
	}

	if err != nil {
		log.Println("Got error when connect database, the error is: " + err.Error())
		return nil, err
	}
	db.LogMode(true)
	if gin.IsDebugging() {
		db.LogMode(true)
	}

	if os.Getenv("AUTOMIGRATE") == "1" {
		db.AutoMigrate(
			&models.User{},
			&models.Invite{},
			&models.Log{},
			&models.RefreshCode{},
		)
	}
	return db, nil
}

func DBInstance(c *gin.Context) *gorm.DB {
	return c.MustGet("DB").(*gorm.DB)
}

func (self *Parameter) SetPreloads(db *gorm.DB) *gorm.DB {
	if self.Preloads == "" {
		return db
	}
	for _, preload := range strings.Split(self.Preloads, ",") {
		var a []string
		for _, s := range strings.Split(preload, ".") {
			a = append(a, snaker.SnakeToCamel(s))
		}
		db = db.Preload(strings.Join(a, "."))
	}
	return db
}
