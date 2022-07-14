package controllers

import (
	dbpkg "github.com/filiponegrao/maromba-back/db"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/gin-gonic/gin"
)

func SaveLog(c *gin.Context, log models.Log) (models.Log, error) {
	db := dbpkg.DBInstance(c)
	if err := db.Save(&log).Error; err != nil {
		return log, err
	} else {
		return log, nil
	}
}

func GetLogs(c *gin.Context) {
	user := GetUserLogged(c)
	db := dbpkg.DBInstance(c)
	// Verifica se eh administrador
	if !user.Admin {
		message := "Apenas usuários administradores possuem acesso!"
		RespondError(c, message, 400)
		return
	}
	var logs []models.Log
	parameter, err := dbpkg.NewParameter(c, models.Log{})
	if err != nil {
		RespondError(c, err.Error(), 400)
		return
	}

	db, err = parameter.Paginate(db)
	if err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	db = parameter.SortRecords(db)

	maxDate := c.Query("maxDate")
	if maxDate != "" {
		db = db.Where("created_at < ?", maxDate)
	}
	searchTerm := c.Query("searchTerm")
	if searchTerm != "" {
		searchQuery := "%" + searchTerm + "%"
		searchQueryPredicate := "ip like ? or url like ? or method like ?"
		searchQueryPredicate += " or user_name like ? or user_mail like ?"
		searchQueryPredicate += " or body like ? or response_body like ?"
		db = db.Where(
			searchQueryPredicate,
			searchQuery,
			searchQuery,
			searchQuery,
			searchQuery,
			searchQuery,
			searchQuery,
			searchQuery,
		)
	}

	if err := db.Find(&logs).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	c.JSON(200, logs)
}
