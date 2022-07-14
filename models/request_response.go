package models

type RequestResponse struct {
	Success bool        `json:"success" form:"success"`
	Result  interface{} `json:"result" form:"result"`
	Message string      `json:"message" form:"message"`
}
