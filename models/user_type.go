package models

// Estrutura nao presente no banco de dados
// Apenas DTO
type UserType struct {
	Type int    `json:"type" form:"type"`
	Name string `json:"name" form:"name"`
}
