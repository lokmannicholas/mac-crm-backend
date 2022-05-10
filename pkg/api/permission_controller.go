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
		_const.PERMISSION_BRANCH.Read(),
		_const.PERMISSION_BRANCH.Create(),
		_const.PERMISSION_BRANCH.Update(),
		_const.PERMISSION_CATEGORY.Read(),
		_const.PERMISSION_CATEGORY.Create(),
		_const.PERMISSION_CATEGORY.Update(),
		_const.PERMISSION_COMPANY.Read(),
		_const.PERMISSION_COMPANY.Create(),
		_const.PERMISSION_COMPANY.Update(),
		_const.PERMISSION_CUSTOMER.Read(),
		_const.PERMISSION_CUSTOMER.Create(),
		_const.PERMISSION_CUSTOMER.Update(),
		_const.PERMISSION_FEATURE.Read(),
		_const.PERMISSION_FEATURE.Create(),
		_const.PERMISSION_FEATURE.Update(),
		_const.PERMISSION_INVOICE.Read(),
		_const.PERMISSION_INVOICE.Create(),
		_const.PERMISSION_INVOICE.Update(),
		_const.PERMISSION_RENTAL_ORDER.Read(),
		_const.PERMISSION_RENTAL_ORDER.Create(),
		_const.PERMISSION_RENTAL_ORDER.Update(),
		_const.PERMISSION_RENTAL_ORDER_EVENT.Read(),
		_const.PERMISSION_RENTAL_ORDER_EVENT.Create(),
		_const.PERMISSION_RENTAL_ORDER_EVENT.Update(),
		_const.PERMISSION_STORAGE.Read(),
		_const.PERMISSION_STORAGE.Create(),
		_const.PERMISSION_STORAGE.Update(),
		_const.PERMISSION_STORAGE_EVENT.Read(),
		_const.PERMISSION_STORAGE_EVENT.Create(),
		_const.PERMISSION_STORAGE_EVENT.Update(),
		_const.PERMISSION_STORAGE_RECORDS.Read(),
		_const.PERMISSION_STORAGE_RECORDS.Create(),
		_const.PERMISSION_STORAGE_RECORDS.Update(),
		_const.PERMISSION_REPORT.Read(),
		_const.PERMISSION_REPORT.Create(),
		_const.PERMISSION_REPORT.Update(),
		_const.PERMISSION_PAYMENT.Read(),
		_const.PERMISSION_PAYMENT.Create(),
		_const.PERMISSION_PAYMENT.Update(),
		_const.PERMISSION_ROLE.Read(),
		_const.PERMISSION_ROLE.Create(),
		_const.PERMISSION_ROLE.Update(),
		_const.PERMISSION_ENCRYPTION.Read(),
		_const.PERMISSION_ENCRYPTION.Create(),
		_const.PERMISSION_ENCRYPTION.Update(),
		_const.PERMISSION_SINGLE_PRODUCT.Read(),
		_const.PERMISSION_SINGLE_PRODUCT.Create(),
		_const.PERMISSION_SINGLE_PRODUCT.Update(),
		_const.PERMISSION_CONSUMABLE_PRODUCT.Read(),
		_const.PERMISSION_CONSUMABLE_PRODUCT.Create(),
		_const.PERMISSION_CONSUMABLE_PRODUCT.Update(),
		_const.PERMISSION_SETTING.Read(),
		_const.PERMISSION_SETTING.Create(),
		_const.PERMISSION_SETTING.Update(),
	}

	controller.Response(c, 200, data)
}
