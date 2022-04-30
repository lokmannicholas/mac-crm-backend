package api

import (
	"fmt"
	"strings"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/lib/auth/middleware"
	"dmglab.com/mac-crm/pkg/util"
	_const "dmglab.com/mac-crm/pkg/util/const"

	"dmglab.com/mac-crm/pkg/entities"
	"dmglab.com/mac-crm/pkg/lib/controller"
	"dmglab.com/mac-crm/pkg/managers"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

type IAccountController interface {
	SetRouter(r *gin.RouterGroup) *AccountController
	GetAccounts(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	GetAccount(c *gin.Context)
	GetAuthAccount(c *gin.Context)
	Logout(c *gin.Context)
	Login(c *gin.Context)
	Disable(c *gin.Context)
	Activate(c *gin.Context)
	ChangePassword(c *gin.Context)
}
type AccountController struct {
	accMgr    managers.IAccountManager
	validator *validator.Validate
}

func NewAccountController() IAccountController {

	return &AccountController{
		accMgr:    managers.GetAccountManager(),
		validator: validator.New(),
	}
}
func (ctl *AccountController) SetRouter(r *gin.RouterGroup) *AccountController {
	accMidWare := middleware.NewAuthMiddleware()

	r.GET("/auth", ctl.GetAuthAccount)
	r.GET("/accounts", accMidWare.PermissionRequire(_const.PERMISSION_ACCOUNT.Read()), ctl.GetAccounts)
	r.POST("/accounts", accMidWare.PermissionRequire(_const.PERMISSION_ACCOUNT.Create()), ctl.Create)
	account := r.Group("/account", accMidWare.AuthRequired())
	{
		account.GET("/:id", accMidWare.PermissionRequire(_const.PERMISSION_ACCOUNT.Read()), ctl.GetAccount)
		account.PUT("/:id", accMidWare.PermissionRequire(_const.PERMISSION_ACCOUNT.Update()), ctl.Update)
		account.PUT("/:id/password/change", accMidWare.PermissionRequire(_const.PERMISSION_ACCOUNT.Update()), ctl.ChangePassword)
		account.POST("/:id/activate", accMidWare.PermissionRequire(_const.PERMISSION_ACCOUNT.Update()), ctl.Activate)
		account.POST("/:id/disable", accMidWare.PermissionRequire(_const.PERMISSION_ACCOUNT.Update()), ctl.Disable)
	}
	return ctl
}

func (ctl *AccountController) GetAccounts(c *gin.Context) {

	accounts, err := ctl.accMgr.GetAccounts(c)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get accounts failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["accounts"] = entities.NewAccountListEntity(c, int64(len(accounts)), accounts)
	controller.Response(c, 200, data)
}

func (ctl *AccountController) Create(c *gin.Context) {

	param := &managers.AccountCreateParam{}
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	err = ctl.validator.Struct(*param)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			controller.ErrorResponse(c, 500, "000000", "account create failed", err.Error())
		} else {
			for _, err := range err.(validator.ValidationErrors) {
				controller.ErrorResponse(c, 500, "000000", "account create failed", fmt.Sprintf(
					"%s has error in validation with %s ", err.Field(), err.Tag()))
				return
			}
		}
		return
	}

	account, err := ctl.accMgr.Create(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "account create failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["account"] = entities.NewAccountEntity(account)
	controller.Response(c, 200, data)
}

func (ctl *AccountController) Update(c *gin.Context) {

	id := c.Param("id")
	param := &managers.AccountUpdateParam{}
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}

	if len(param.RoleID) == 0 || len(param.Status) == 0 {
		controller.ErrorResponse(c, 500, "000000", "account update failed", "Status or Role cannot be empty")
		return
	}

	account, err := ctl.accMgr.Update(c, id, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "account update failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["account"] = entities.NewAccountEntity(account)
	controller.Response(c, 200, data)
}

func (ctl *AccountController) Login(c *gin.Context) {

	param := &managers.LoginParam{}
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	account, err := ctl.accMgr.Login(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "login failed", err.Error())
		return
	}
	c.Set("Account", account)
	token, refreshToken, err := ctl.accMgr.CreateToken(c)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "login failed", err.Error())
		return
	}
	permissions := map[string]interface{}{}
	for _, p := range strings.Split(account.Role.Permissions, ";") {
		permissions[p] = true
	}
	data := map[string]interface{}{
		"authUser": map[string]interface{}{
			"companyID": config.GetConfig().CompanyID,
			"username":  account.Username,
			"role": map[string]interface{}{
				"id":          account.Role.ID,
				"permissions": permissions,
			},
		},
		"token": map[string]string{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
	}
	controller.Response(c, 200, data)
}

func (ctl *AccountController) GetAuthAccount(c *gin.Context) {

	account, err := util.GetCtxAccount(c)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get account failed", err.Error())
		return
	}
	account, err = ctl.accMgr.GetAccount(c, account.ID.String())
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get account failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["account"] = entities.NewAccountEntity(account)
	controller.Response(c, 200, data)
}

func (ctl *AccountController) GetAccount(c *gin.Context) {

	id := c.Param("id")
	account, err := ctl.accMgr.GetAccount(c, id)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get account failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["account"] = entities.NewAccountEntity(account)
	controller.Response(c, 200, data)
}

func (ctl *AccountController) Logout(c *gin.Context) {

	err := ctl.accMgr.Logout(c)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "logout failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	controller.Response(c, 200, data)
}

func (ctl *AccountController) ChangePassword(c *gin.Context) {

	id := c.Param("id")
	param := &managers.ChangePasswordParam{}
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}

	err = ctl.accMgr.ChangePassword(c, id, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "change password failed", err.Error())
		return
	}
	token, refreshToken, err := ctl.accMgr.CreateToken(c)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "login failed", err.Error())
		return
	}
	data := map[string]interface{}{
		"token": map[string]string{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
	}
	controller.Response(c, 200, data)
}

func (ctl *AccountController) Activate(c *gin.Context) {

	id := c.Param("id")

	err := ctl.accMgr.Activate(c, id)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "account modify failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	controller.Response(c, 200, data)
}

func (ctl *AccountController) Disable(c *gin.Context) {

	id := c.Param("id")

	err := ctl.accMgr.Disable(c, id)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "account modify failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	controller.Response(c, 200, data)
}
