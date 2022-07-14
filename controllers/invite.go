package controllers

import (
	"time"

	dbpkg "github.com/filiponegrao/maromba-back/db"
	"github.com/filiponegrao/maromba-back/email"
	"github.com/filiponegrao/maromba-back/models"
	"github.com/filiponegrao/maromba-back/tools"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/************************************************
/**** MARK: AUXILIARY METHODS ****/
/************************************************/
func CreateInvite(
	c *gin.Context, db *gorm.DB,
	code string, user models.User,
	password string,
) (invite models.Invite, err error) {
	inviter := GetUserLogged(c)
	invite.InviterID = inviter.ID
	invite.InvitedID = user.ID
	invite.Code = tools.EncryptTextSHA512(code)
	invite.Status = models.INVITE_STATUS_PENDING
	invite.Inviter = inviter
	invite.Invited = user
	// Cria o convite
	if err := db.Create(&invite).Error; err != nil {
		return invite, err
	}
	// Envia o email
	if user.Type == 0 {
		return invite, email.UserPendingInvite(Conf, invite, code)
	} else {
		return invite, email.CollaboratorInvite(Conf, invite, password, code)
	}
}

/* Valida um cadastro convidado */
func ValidateInvite(c *gin.Context) {
	user := GetUserLogged(c)
	db := dbpkg.DBInstance(c)
	code := c.Params.ByName("code")
	encodedCode := tools.EncryptTextSHA512(code)
	var invite models.Invite
	if err := db.Where(
		"status = ? AND code = ? AND invited_id = ?",
		models.INVITE_STATUS_PENDING,
		encodedCode, user.ID,
	).First(&invite).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	tx := db.Begin()
	// Altera o status do usuario
	if user.Status == models.USER_STATUS_PENDING {
		user.Status = models.USER_STATUS_AVAILABLE
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			RespondError(c, err.Error(), 400)
			return
		}
	}
	now := time.Now()
	// Invalida o invite
	invite.Status = models.INVITE_STATUS_VALIDATED
	invite.UpdatedAt = &now
	if err := tx.Save(&invite).Error; err != nil {
		tx.Rollback()
		RespondError(c, err.Error(), 400)
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		RespondError(c, err.Error(), 400)
		return
	}
	RespondSuccessV2(c, nil)
}

func ResendInvite(c *gin.Context) {
	user := GetUserLogged(c)
	db := dbpkg.DBInstance(c)
	var invite models.Invite
	if err := db.Where(
		"status = ? AND invited_id = ?",
		models.INVITE_STATUS_PENDING,
		user.ID,
	).First(&invite).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	if err := db.First(&invite.Invited, invite.InvitedID).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	db.First(&invite.Inviter, invite.InviterID)

	// Cria um novo código
	code := tools.RandomNumbers(6)
	invite.Code = tools.EncryptTextSHA512(code)
	tx := db.Begin()
	if err := tx.Save(&invite).Error; err != nil {
		tx.Rollback()
		RespondError(c, err.Error(), 400)
		return
	}
	// Envia o email
	if invite.Invited.Type == 0 {
		if err := email.UserPendingInvite(Conf, invite, code); err != nil {
			tx.Rollback()
			RespondError(c, err.Error(), 400)
			return
		}
	} else {
		if err := email.ReSendInvite(Conf, invite, code); err != nil {
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
	RespondSuccess(c, nil)
}

/* Lista todos os convites enviados */
func GetSentInvites(c *gin.Context) {
	user := GetUserLogged(c)
	db := dbpkg.DBInstance(c)
	if user.Type != models.USER_TYPE_ADMIN &&
		user.Type != models.USER_TYPE_MANAGER {
		message := "Sem acesso!"
		RespondError(c, message, 400)
		return
	}
	var invites []models.Invite
	if err := db.Where("inviter_id = ?", user.ID).Find(&invites).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	for i := 0; i < len(invites); i++ {
		db.First(&invites[i].Invited, invites[i].InvitedID)
		invites[i].Invited.Password = ""
	}
	RespondSuccessV2(c, invites)
}

/* Efetua a troca de perfil de um colaborador (gerente ou não) */
func GetCollaboratorInvite(c *gin.Context) {
	user := GetUserLogged(c)
	if user.Type != models.USER_TYPE_ADMIN &&
		user.Type != models.USER_TYPE_MANAGER {
		message := "Sem acesso!"
		RespondError(c, message, 400)
		return
	}
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	var invite models.Invite
	if err := db.Where("invited_id = ?", id).First(&invite).Error; err != nil {
		RespondError(c, err.Error(), 400)
		return
	}
	db.First(&invite.Inviter, invite.InviterID)
	db.First(&invite.Invited, invite.InvitedID)
	invite.Inviter.Password = ""
	invite.Invited.Password = ""
	RespondSuccessV2(c, invite)
}
