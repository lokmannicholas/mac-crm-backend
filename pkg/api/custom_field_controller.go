package api

import (
	"dmglab.com/mac-crm/pkg/entities"
	"dmglab.com/mac-crm/pkg/lib/auth/middleware"
	controller "dmglab.com/mac-crm/pkg/lib/controller"
	"dmglab.com/mac-crm/pkg/managers"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/gin-gonic/gin"
)

type ICustomFieldController interface {
	SetRouter(r *gin.RouterGroup) *CustomFieldController
	GetCustomFields(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
}
type CustomFieldController struct {
	cusMgr managers.ICustomFieldManager
}

func NewCustomFieldController() ICustomFieldController {
	return &CustomFieldController{
		cusMgr: managers.GetCustomFieldManager(),
	}
}

func (ctl *CustomFieldController) SetRouter(r *gin.RouterGroup) *CustomFieldController {
	accMidWare := middleware.NewAuthMiddleware()
	r.GET("/custom-fields", accMidWare.PermissionRequire(_const.PERMISSION_SETTING.Read()), ctl.GetCustomFields)
	r.POST("/custom-fields", accMidWare.PermissionRequire(_const.PERMISSION_SETTING.Create()), ctl.Create)
	cus := r.Group("/custom-field")
	{
		cus.PUT("/:id", accMidWare.PermissionRequire(_const.PERMISSION_SETTING.Update()), ctl.Update)

	}
	return ctl
}

// GetCustomFields godoc
// @Tags Custom Field
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param custom_object query string false " "
// @Param status query string false " "
// @Param page query int false " "
// @Param limit query int false " "
// @Success 200 {object} swagger.APIResponse{data=swagger.CustomFields}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /custom-fields [get]
func (ctl *CustomFieldController) GetCustomFields(c *gin.Context) {

	param := new(managers.CustomFieldQueryParam)
	err := controller.GetQueryParam(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get CustomFields failed", err.Error())
		return
	}
	CustomFields, pagin, err := ctl.cusMgr.GetCustomFields(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get CustomFields failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["custom_fields"] = entities.NewCustomFieldListEntity(pagin.TotalCount, CustomFields)

	controller.Response(c, 200, data)
}

// CreateCustomFields godoc
// @Tags Custom Field
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param payload body managers.CustomFieldCreateParam true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.CustomField}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /custom-fields [post]
func (ctl *CustomFieldController) Create(c *gin.Context) {

	param := new(managers.CustomFieldCreateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	CustomField, err := ctl.cusMgr.Create(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "create CustomFields failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["custom_field"] = entities.NewCustomFieldEntity(CustomField)
	controller.Response(c, 200, data)
}

// UpdateCustomField godoc
// @Tags Custom Field
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param id path string true " "
// @Param payload body managers.CustomFieldUpdateParam true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.CustomField}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /custom-field/:id [put]
func (ctl *CustomFieldController) Update(c *gin.Context) {

	CustomFieldID := c.Param("id")
	param := new(managers.CustomFieldUpdateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	CustomField, err := ctl.cusMgr.Update(c, CustomFieldID, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "create CustomFields failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["custom_field"] = entities.NewCustomFieldEntity(CustomField)
	controller.Response(c, 200, data)
}
