package models

import "time"

type Log struct {
	ID           int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	IP           string     `gorm:"ip" form:"ip" json:"ip"`
	UserMail     string     `gorm:"default:''" json:"userMail" form:"userMail"`
	UserID       int64      `gorm:"default:0" json:"userId" form:"userId"`
	UserName     string     `gorm:"default:''" json:"userName" form:"userName"`
	Method       string     `gorm:"not null" json:"method" form:"method"`
	URL          string     `gorm:"not null" json:"url" form:"url"`
	Body         string     `gorm:"default:''" json:"body" form:"body"`
	Header       string     `gorm:"default:''" json:"header" form:"header"`
	Description  string     `gorm:"default:''" json:"description" form:"description"`
	ResponseCode int        `gorm:"not null" json:"responseCode" form:"responseCode"`
	ResponseBody string     `gorm:"default:''" json:"responseBody" form:"responseBody"`
	CreatedAt    *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" form:"updated_at"`
}

func (log Log) GetDescriptionString() string {
	// verb := ""
	// if log.Method == "GET" {
	// 	verb = "Visualização de "
	// } else if log.Method == "POST" {
	// 	verb = "Registro de "
	// } else if log.Method == "PUT" {
	// 	verb = "Alteração de "
	// } else if log.Method == "DELETE" {
	// 	verb = "Exclusão de "
	// }
	return ""
}
