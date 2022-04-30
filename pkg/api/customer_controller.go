package api

import (
	"dmglab.com/mac-crm/pkg/entities"
	controller "dmglab.com/mac-crm/pkg/lib/controller"
	"dmglab.com/mac-crm/pkg/managers"
	"github.com/gin-gonic/gin"
)

type ICustomerController interface {
	GetCustomers(c *gin.Context)
	GetCustomer(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
}
type CustomerController struct {
	cusMgr managers.ICustomerManager
}

func NewCustomerController() ICustomerController {
	return &CustomerController{
		cusMgr: managers.GetCustomerManager(),
	}
}

func (ctl *CustomerController) GetCustomers(c *gin.Context) {

	param := new(managers.CustomerQueryParam)
	if err := controller.GetQueryParam(c, param); err != nil {
		controller.ErrorResponse(c, 500, "000000", "get customers failed", err.Error())
		return
	}
	customers, pagin, err := ctl.cusMgr.GetCustomers(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get customers failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["customers"] = entities.NewCustomerListEntity(pagin.TotalCount, customers)

	controller.Response(c, 200, data)
}

func (ctl *CustomerController) GetCustomer(c *gin.Context) {

	customerID := c.Param("id")
	var err error
	customer, err := ctl.cusMgr.GetCustomer(c, customerID)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get customer failed", err.Error())
		return
	}
	data := map[string]interface{}{}

	data["customer"] = entities.NewCustomerEntity(customer)
	controller.Response(c, 200, data)
}

func (ctl *CustomerController) Create(c *gin.Context) {

	param := new(managers.CustomerCreateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	customer, err := ctl.cusMgr.Create(c, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "create customers failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["customer"] = entities.NewCustomerEntity(customer)
	controller.Response(c, 200, data)
}

func (ctl *CustomerController) Update(c *gin.Context) {

	customerID := c.Param("id")
	param := new(managers.CustomerUpdateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	customer, err := ctl.cusMgr.Update(c, customerID, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "create customers failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["customer"] = entities.NewCustomerEntity(customer)
	controller.Response(c, 200, data)
}
