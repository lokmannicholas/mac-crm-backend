package api

import (
	"dmglab.com/mac-crm/pkg/entities"
	"dmglab.com/mac-crm/pkg/lib/auth/middleware"
	controller "dmglab.com/mac-crm/pkg/lib/controller"
	"dmglab.com/mac-crm/pkg/managers"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ISettingController interface {
	SetRouter(r *gin.RouterGroup) *SettingController
	Set(c *gin.Context)
	Get(c *gin.Context)
}
type SettingController struct {
	settingMgr managers.ISettingManager
}

func NewSettingController() ISettingController {
	return &SettingController{
		settingMgr: managers.GetSettingManager(),
	}
}
func (ctl *SettingController) SetRouter(r *gin.RouterGroup) *SettingController {
	accMidWare := middleware.NewAuthMiddleware()
	r.GET("/setting/general", accMidWare.PermissionRequire(_const.PERMISSION_SETTING.Read()), ctl.Get)
	r.POST("/setting/general", accMidWare.PermissionRequire(_const.PERMISSION_SETTING.Create()), ctl.Set)

	return ctl
}

// SetSetting godoc
// @Tags Setting
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param payload body managers.SettingSetParam true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.Settings}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /setting/general [post]
func (ctl *SettingController) Set(c *gin.Context) {
	param := &managers.SettingSetParam{}
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}

	appSetting, err := ctl.settingMgr.Set(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "create Setting failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["settings"] = entities.NewSettingListEntity(int64(len(appSetting)), appSetting)
	controller.Response(c, 200, data)
}

// GetSettings godoc
// @Tags Setting
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.Settings}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /setting/general [get]
func (ctl *SettingController) Get(c *gin.Context) {

	param := &managers.SettingGetParam{}
	err := controller.GetQueryParam(c, param)
	if err != nil {
		return
	}

	data := map[string]interface{}{}
	if param.Key != nil {
		setting, err := ctl.settingMgr.GetByKey(c, param)
		if err != nil {
			controller.ErrorResponse(c, 500, "000000", "get Setting failed", err.Error())
			return
		}
		data["setting"] = entities.NewSettingEntity(setting)
	} else {
		db, ok := c.Get("DB")
		if !ok {
			controller.ErrorResponse(c, 400, "000000", "get Setting failed", "")
			return
		}
		settings, err := ctl.settingMgr.GetAll(db.(*gorm.DB))
		if err != nil {
			controller.ErrorResponse(c, 500, "000000", "get Setting failed", err.Error())
			return
		}
		data["settings"] = entities.NewSettingListEntity(int64(len(settings)), settings)
	}

	controller.Response(c, 200, data)
}
