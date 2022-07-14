package controllers

import (
	"strings"

	dbpkg "github.com/filiponegrao/maromba-back/db"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/gin-gonic/gin"
)

/************************************************
/**** MARK: AUXILIARY METHODS ****/
/************************************************/

func CheckUserExists(c *gin.Context, email string) (exists bool, err error, user models.User) {
	db := dbpkg.DBInstance(c)
	var existent models.User
	if err = db.Where("email = ?", email).First(&existent).Error; err != nil {
		errorMessage := err.Error()
		if strings.Contains(errorMessage, "record not found") {
			return false, nil, existent
		} else {
			return false, err, existent
		}
	}
	if existent.ID != 0 {
		return true, nil, existent
	} else {
		return false, nil, existent
	}
}

func GetTotalUsersCount() (usersCount int, err error) {
	return usersCount, nil
}

func GetNewUsersCount() {
}

func GetTotalPendingUsers(usersCount int, err error) {
}
