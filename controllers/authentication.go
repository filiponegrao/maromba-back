package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"

	dbpkg "github.com/filiponegrao/maromba-back/db"
	"github.com/filiponegrao/maromba-back/email"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/filiponegrao/maromba-back/tools"
	"github.com/gin-gonic/gin"
)

// MARK: Authentication

const TimeOutToken = time.Minute * 30
const MaxRefreshToken = time.Hour * 24 * 30

// GetAuthMiddlware : GetAuthMiddlware
func GetAuthMiddlware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone 2",
		Key:             []byte("A@c#e$ler&ad*osS(erv)er01"),
		Timeout:         TimeOutToken,
		MaxRefresh:      MaxRefreshToken,
		IdentityKey:     "id",
		PayloadFunc:     AuthorizationPayload,
		IdentityHandler: IdentityHandler,
		Authenticator:   UserAuthentication,
		Authorizator:    UserAuthorization,
		Unauthorized:    UserUnauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
		LoginResponse:   LoginResponse,
		RefreshResponse: LoginResponse,
	})
	return authMiddleware, err
}

// MARK: Structures

type objectLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Token    string `form:"token" json:"token"`
}

type objectNewPassword struct {
	OldPassowrd     string `json:"old_password" form:"old_password"`
	NewPassword     string `json:"new_password" form:"new_password"`
	ConfirmPassowrd string `json:"confirm_password" form:"confirm_password"`
}

type loginResponse struct {
	Token       string     `json:"token" form:"token"`
	RefreshCode string     `json:"refresh_code" form:"refresh_code"`
	Expire      *time.Time `json:"expire" form:"expire"`
}

// ChangePassword : ChangePassword
func ChangePassword(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	var encoded string
	var newPassword objectNewPassword

	if err := c.Bind(&encoded); err != nil {

		RespondError(c, err.Error(), 400)
		return
	}

	sDec, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	if err := json.Unmarshal([]byte(sDec), &newPassword); err != nil {
		RespondError(c, err.Error(), 400)
		return
	}

	if newPassword.OldPassowrd == "" {
		message := "Faltando senha atual."
		RespondError(c, message, 400)
		return
	} else if newPassword.NewPassword == "" {
		message := "Faltando nova senha."
		RespondError(c, message, 400)
		return
	} else if newPassword.ConfirmPassowrd == "" {
		message := "Faltando confirmacao da nova senha."
		RespondError(c, message, 400)
		return
	}

	claims := jwt.ExtractClaims(c)
	userID := int64(claims["id"].(float64))

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	encPassword := GetAuthorizationToken(user.Email, newPassword.OldPassowrd)
	// Vrifica corretude de senha recebida
	if user.Password != encPassword {
		message := "Senha atual incorreta."
		RespondError(c, message, 400)
		return
	}
	// Verifica se nova senha foi escrita corretamente
	if newPassword.NewPassword != newPassword.ConfirmPassowrd {
		message := "Nova senha nao confere."
		RespondError(c, message, 400)
		return
	}

	newPasswordEnc := GetAuthorizationToken(user.Email, newPassword.NewPassword)
	user.Password = newPasswordEnc
	if err := db.Save(user).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}

	// EmailChangedPassword(user.Email)
	RespondSuccess(c, "Senha atualizada com sucesso!")
}

// ForgotPassword : Forgot Password
func ForgotPassword(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	body := c.Request.RequestURI
	var user models.User

	parts := strings.Split(body, "email=")
	if len(parts) <= 1 {
		RespondError(c, "Faltando parametro de email", 400)
		return
	}

	emailText := parts[1]
	if strings.Contains(emailText, "'") {
		message := "E-mail incorreto"
		RespondError(c, message, 400)
		return
	}
	if err := db.Where("email = ?", emailText).First(&user).Error; err != nil {
		errorMEssage := err.Error()
		if errorMEssage != "record not found" {
			RespondError(c, errorMEssage, 400)
			return
		} else {
			RespondError(c, "Usuário não encontrado!", 400)
			return
		}
	}

	password := tools.RandomString(6)
	passwordEncode := GetAuthorizationToken(emailText, password)
	user.Password = passwordEncode

	if err := db.Save(user).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}

	email.EmailPasswordNew(Conf, user.Email, password)

	RespondSuccess(c, "Nova senha enviada para o email!")
}

// UserAuthentication : UserAuthentication
func UserAuthentication(c *gin.Context) (interface{}, error) {
	var loginVals objectLogin
	if err := c.Bind(&loginVals); err != nil {
		return nil, err
	}
	email := loginVals.Username
	password := loginVals.Token

	db := dbpkg.DBInstance(c)
	if email == "" {
		message := "Faltando email"
		return nil, errors.New(message)
	}
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("Usuário não cadastrado")
	}

	if password != user.Password {
		SaveLogData(c, 401, "Senha incorreta")
		return nil, errors.New("Senha incorreta")
	}

	user.Password = ""
	return &user, nil
}

// UserAuthorization : UserAuthorization
func UserAuthorization(user interface{}, c *gin.Context) bool {
	return true
}

// UserUnauthorized : Falha na autênticação
func UserUnauthorized(c *gin.Context, code int, message string) {
	err := ""
	if strings.Contains(message, "missing") {
		err = "Faltando email ou senha"
	} else if strings.Contains(message, "incorrect") {
		err = "Email ou senha incorreta"
	} else if strings.Contains(message, "cookie token is empty") {
		err = "Faltando HEADER de autenticação!"
		// } else if strings.Contains(message, "") {

	} else {
		err = message
	}
	RespondError(c, err, 400)
}

// AuthorizationPayload : AuthorizationPayload
func AuthorizationPayload(data interface{}) jwt.MapClaims {
	if user, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			"id": user.ID,
		}
	}
	return jwt.MapClaims{}
}

// IdentityHandler : IdentityHandler
func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &models.User{
		ID: int64(claims["id"].(float64)),
	}
}

// GetAuthorizationToken : GetAuthorizationToken
func GetAuthorizationToken(email string, password string) string {
	encPassword := tools.EncryptTextSHA512(password)
	encPassword = email + ":" + encPassword
	encPassword = tools.EncryptTextSHA512(encPassword)
	return encPassword
}
