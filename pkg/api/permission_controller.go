package api

import (
	"dmglab.com/mac-crm/pkg/lib/controller"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/gin-gonic/gin"
)

type IPermissionController interface {
	GetPermissions(c *gin.Context)
}
type PermissionController struct {
}

func NewPermissionController() IPermissionController {
	return &PermissionController{}
}

// GetPermissions godoc
// @Tags Permission
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.Permissions}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /permissions [get]
func (ctl *PermissionController) GetPermissions(c *gin.Context) {
	data := map[string]interface{}{}
	data["permissions"] = []string{
		_const.PERMISSION_ACCOUNT.Read(),
		_const.PERMISSION_ACCOUNT.Create(),
		_const.PERMISSION_ACCOUNT.Update(),
		_const.PERMISSION_APP.Read(),
		_const.PERMISSION_APP.Create(),
		_const.PERMISSION_APP.Update(),
		_const.PERMISSION_COMPANY.Read(),
		_const.PERMISSION_COMPANY.Create(),
		_const.PERMISSION_COMPANY.Update(),
		_const.PERMISSION_CUSTOMER.Read(),
		_const.PERMISSION_CUSTOMER.Create(),
		_const.PERMISSION_CUSTOMER.Update(),
		_const.PERMISSION_ROLE.Read(),
		_const.PERMISSION_ROLE.Create(),
		_const.PERMISSION_ROLE.Update(),
		_const.PERMISSION_SETTING.Read(),
		_const.PERMISSION_SETTING.Create(),
		_const.PERMISSION_SETTING.Update(),
	}

	controller.Response(c, 200, data)
}
