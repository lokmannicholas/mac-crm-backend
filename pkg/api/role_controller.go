package api

import (
	"dmglab.com/mac-crm/pkg/entities"
	controller "dmglab.com/mac-crm/pkg/lib/controller"
	"dmglab.com/mac-crm/pkg/managers"
	"github.com/gin-gonic/gin"
)

type IRoleController interface {
	GetRoles(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
}
type RoleController struct {
	rolMgr managers.IRoleManager
}

func NewRoleController() IRoleController {
	return &RoleController{
		rolMgr: managers.GetRoleManager(),
	}
}

func (ctl *RoleController) GetRoles(c *gin.Context) {

	roles, err := ctl.rolMgr.GetRoles(c)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get roles failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["roles"] = entities.NewRoleListEntity(int64(len(roles)), roles)
	controller.Response(c, 200, data)
}

func (ctl *RoleController) Update(c *gin.Context) {

	id := c.Param("id")
	param := new(managers.RoleUpdateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "role update failed", err.Error())
		return
	}

	if len(id) == 0 {
		controller.ErrorResponse(c, 500, "000000", "role update failed", "id is missing")
		return
	}
	if len(param.Name) == 0 {
		controller.ErrorResponse(c, 500, "000000", "role update failed", "name is missing")
		return
	}
	role, err := ctl.rolMgr.Update(c, id, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "role update failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["role"] = entities.NewRoleEntity(role)
	controller.Response(c, 200, data)
}

func (ctl *RoleController) Create(c *gin.Context) {

	param := new(managers.RoleCreateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}

	if len(param.Name) == 0 {
		controller.ErrorResponse(c, 500, "000000", "role create failed", "name is missing")
		return
	}
	roles, err := ctl.rolMgr.Create(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "role create failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["role"] = entities.NewRoleEntity(roles)
	controller.Response(c, 200, data)
}
