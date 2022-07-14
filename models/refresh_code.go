package models

import "time"

type RefreshCode struct {
	ID          int64      `gorm:"primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	AccessToken string     `json:"access_token" form:"access_token"`
	Code        string     `json:"code" form:"cpde"`
	Status      int64      `json:"status" form:"status"`
	CreatedAt   *time.Time `json:"created_at" form:"created_at"`
}
