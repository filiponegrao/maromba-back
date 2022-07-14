package controllers

import (
	"strings"

	dbpkg "github.com/filiponegrao/maromba-back/db"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/filiponegrao/maromba-back/tools"

	"github.com/gin-gonic/gin"
)

/************************************************
/**** MARK: NORMAL USERS ****/
/************************************************/

func CreateUser(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	missing := user.MissingFields()
	if missing != "" {
		message := "Faltando campo " + missing
		RespondError(c, message, 400)
		return
	}
	// Valida e-mail
	if !tools.ValidateEmail(user.Email) {
		message := "E-mail inválido!"
		RespondError(c, message, 400)
		return
	}
	// Verifica se o usuario existe
	exists, err, _ := CheckUserExists(c, user.Email)
	if err != nil {
		RespondError(c, err.Error(), 400)
		return
	} else if exists {
		message := "Usuário já existe"
		RespondError(c, message, 400)
		return
	}
	// Atualiza senha
	if user.Password != "" {
		passwordEncode := tools.EncryptTextSHA512(user.Password)
		passwordEncode = user.Email + ":" + passwordEncode
		passwordEncode = tools.EncryptTextSHA512(passwordEncode)
		user.Password = passwordEncode
	}
	user.Admin = false
	user.Type = models.USER_TYPE_NORMAL
	user.Status = models.USER_STATUS_AVAILABLE
	headerVersion := c.Request.Header.Get("Application-Version")
	if headerVersion == "v1" {
		user.Status = models.USER_STATUS_PENDING
	}
	tx := db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		RespondError(c, err.Error(), 400)
		return
	}
	// Verifica se há necessidade de confirmação de cadastro
	if headerVersion == "v1" {
		user.Status = models.USER_STATUS_PENDING
		code := tools.RandomString(6)
		_, err := CreateInvite(c, tx, code, user, "")
		if err != nil {
			tx.Rollback()
			RespondError(c, err.Error(), 400)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		RespondError(c, err.Error(), 400)
		return
	}

	// Campos sigilosos
	user.Password = ""
	RespondSuccess(c, user)
}

/************************************************
/**** MARK: ALL USERS ****/
/************************************************/

func UpdateUser(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	user := GetUserLogged(c)
	// Parametros que nao irao mudar
	password := user.Password
	admin := user.Admin
	userType := user.Type
	status := user.Status
	imageUrl := user.ProfileImageURL
	cnpj := user.CNPJ
	cpf := user.CPF
	email := user.Email
	if err := c.Bind(&user); err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	// Atualiza somente se nunca houver sido atribuido
	if cpf == "" && user.CPF != "" {
		// Verifica se ha disponibilidade
		existent, err := CheckCPF(c, user.CPF, user.ID)
		if err != nil {
			RespondError(c, err.Error(), 400)
			return
		}
		if existent {
			message := "CPF já reigstrado"
			RespondError(c, message, 400)
			return
		}
	} else {
		user.CPF = cpf
	}
	// Atualiza somente se nunca houver sido atribuido
	if cnpj == "" && user.CNPJ != "" {
		// Verifica se ha disponibilidade
		// existent, err := CheckCNPJ(c, user.CNPJ, user.ID)
		// if err != nil {
		// 	RespondError(c, err.Error(), 400)
		// 	return
		// }
		// if existent {
		// 	message := "CNPJ já reigstrado"
		// 	RespondError(c, message, 400)
		// 	return
		// }
	} else {
		user.CNPJ = cnpj
	}
	user.Password = password
	user.Admin = admin
	user.Type = userType
	user.Status = status
	user.Email = email
	if user.ProfileImageURL == "" && imageUrl != "" {
		user.ProfileImageURL = imageUrl
	}
	user.Phone1 = tools.RemoveSpecialCaracters(user.Phone1)
	user.Phone2 = tools.RemoveSpecialCaracters(user.Phone2)
	// Valida celular
	if user.Phone1 != "" && len(user.Phone1) < 11 {
		RespondError(c, "Número de celular incorreto!", 400)
		return
	}
	// Valida telefone
	if user.Phone2 != "" && len(user.Phone2) < 10 {
		RespondError(c, "Número de telefone incorreto!", 400)
		return
	}
	if err := db.Save(&user).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	// Campos sigilosos
	user.Password = ""
	RespondSuccess(c, user)
}

/* Visualiza informações do usuario logado */
func RequestUserLogged(c *gin.Context) {
	user := GetUserLogged(c)
	user.Password = ""
	RespondSuccess(c, user)
}

// CheckCPF : CheckCPF
func CheckCPF(c *gin.Context, cpf string, user int64) (existent bool, err error) {
	db := dbpkg.DBInstance(c)
	var wanted models.User
	if err := db.Where("cpf = ? AND id != ?", cpf, user).First(&wanted).Error; err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			return false, err
		}
		return false, nil
	}
	return wanted.ID != 0, nil
}

// CheckCNPJ : CheckCNPJ
func CheckCNPJ(c *gin.Context, cnpj string, user int64) (existent bool, err error) {
	db := dbpkg.DBInstance(c)
	var wanted models.User
	if err := db.Where("cnpj = ? AND id != ?", cnpj, user).First(&wanted).Error; err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			return false, err
		}
		return false, nil
	}
	return wanted.ID != 0, nil
}

/************************************************
/**** MARK: ADMIN ****/
/************************************************/

func GetUsers(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	user := GetUserLogged(c)
	if user.Type != models.USER_TYPE_ADMIN {
		message := "Sem acesso!"
		RespondError(c, message, 400)
		return
	}
	// Pagination
	parameter, err := dbpkg.NewParameter(c, models.User{})
	if err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	db, err = parameter.Paginate(db)
	if err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	db = UserSearchTermQueryDb(c, db)
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	for i := 0; i < len(users); i++ {
		users[i].Password = ""
	}
	RespondSuccessV2(c, users)
}

// MARK: Aux

func AddNonRepetitiveUsersToArray(itens []models.User, elements []models.User) []models.User {
	var result []models.User
	result = append(result, itens...)
	for _, element := range elements {
		alreadyExists := false
		for _, item := range itens {
			if item.ID == element.ID {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			result = append(result, element)
		}
	}
	return result
}

func AddNonRepetitiveUserToArray(itens []models.User, element models.User) []models.User {
	var result []models.User
	result = append(result, itens...)
	alreadyExists := false
	for _, item := range itens {
		if item.ID == element.ID {
			alreadyExists = true
			break
		}
	}
	if !alreadyExists {
		result = append(result, element)
	}
	return result
}
