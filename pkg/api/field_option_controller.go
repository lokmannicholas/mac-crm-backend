package api

import (
	"dmglab.com/mac-crm/pkg/entities"
	"dmglab.com/mac-crm/pkg/lib/auth/middleware"
	controller "dmglab.com/mac-crm/pkg/lib/controller"
	"dmglab.com/mac-crm/pkg/managers"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/gin-gonic/gin"
)

type IFieldOptionController interface {
	SetRouter(r *gin.RouterGroup) *FieldOptionController
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
type FieldOptionController struct {
	fieldOptMgr managers.IFieldOptionManager
}

func NewFieldOptionController() IFieldOptionController {
	return &FieldOptionController{
		fieldOptMgr: managers.GetFieldOptionManager(),
	}
}

func (ctl *FieldOptionController) SetRouter(r *gin.RouterGroup) *FieldOptionController {
	accMidWare := middleware.NewAuthMiddleware()
	cus := r.Group("/field-option")
	{
		cus.POST("/", accMidWare.PermissionRequire(_const.PERMISSION_SETTING.Create()), ctl.Create)
		cus.PUT("/:id", accMidWare.PermissionRequire(_const.PERMISSION_SETTING.Update()), ctl.Update)
		cus.DELETE("/:id", accMidWare.PermissionRequire(_const.PERMISSION_SETTING.Update()), ctl.Delete)
	}
	return ctl
}

// CreateFieldOption godoc
// @Tags Field Option
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param payload body managers.FieldOptionCreateParam true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.FieldOption}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /field-option [post]
func (ctl *FieldOptionController) Create(c *gin.Context) {

	param := new(managers.FieldOptionCreateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	FieldOption, err := ctl.fieldOptMgr.Create(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "create FieldOption failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["field_option"] = entities.NewFieldOptinoEntity(FieldOption)
	controller.Response(c, 200, data)
}

// UpdateFieldOption godoc
// @Tags Field Option
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param id path string true " "
// @Param payload body managers.FieldOptionUpdateParam true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.FieldOption}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /field-option/:id [put]
func (ctl *FieldOptionController) Update(c *gin.Context) {

	OptionID := c.Param("id")
	param := new(managers.FieldOptionUpdateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	FieldOption, err := ctl.fieldOptMgr.Update(c, OptionID, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "update FieldOption failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["field_option"] = entities.NewFieldOptinoEntity(FieldOption)
	controller.Response(c, 200, data)
}

// DeleteFieldOption godoc
// @Tags Field Option
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param id path string true " "
// @Success 200
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /field-option/:id [delete]
func (ctl *FieldOptionController) Delete(c *gin.Context) {

	OptionID := c.Param("id")
	err := ctl.fieldOptMgr.Delete(c, OptionID)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "delete FieldOption failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	controller.Response(c, 200, data)
}
