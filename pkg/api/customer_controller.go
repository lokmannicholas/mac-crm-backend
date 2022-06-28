package api

import (
	"dmglab.com/mac-crm/pkg/entities"
	controller "dmglab.com/mac-crm/pkg/lib/controller"
	"dmglab.com/mac-crm/pkg/managers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// GetCustomers godoc
// @Tags Customer
// @Description meta value will be string or entities.FieldOption array
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param first_name query string false " "
// @Param last_name query string false " "
// @Param id_no query string false " "
// @Param court_order_date query string false " " example(2022-05-14T00:00:00.000Z#2022-07-14T00:00:00.000Z)
// @Param court_release_date query string false " " example(2022-05-14T00:00:00.000Z#2022-07-14T00:00:00.000Z)
// @Param email query string false " "
// @Param phone1 query string false " "
// @Param phone2 query string false " "
// @Param phone3 query string false " "
// @Param loan_type query string false " "
// @Param is_bankrupt query string false " "
// @Param is_drp query string false " "
// @Param is_iva query string false " "
// @Param court_case_involved query string false " "
// @Param department query string false " "
// @Success 200 {object} swagger.APIResponse{data=swagger.Customers}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /customers [get]
func (ctl *CustomerController) GetCustomers(c *gin.Context) {

	param := new(managers.CustomerQueryParam)
	if err := controller.GetQueryParam(c, param); err != nil {
		controller.ErrorResponse(c, 500, "000000", "get customers failed", err.Error())
		return
	}

	v, ok := c.Get("FieldPermissions")
	if !ok {
		c.AbortWithStatusJSON(403, gin.H{"message": "Permission denied"})
		return
	}
	fieldPermissions := v.(string)

	levels := ""
	v, ok = c.Get("Levels")
	if ok {
		levels = v.(string)
	}

	customers, err := ctl.cusMgr.GetCustomers(c, param, fieldPermissions, levels)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get customer failed", err.Error())
		return
	}
	data := map[string]interface{}{}

	data["customers"] = entities.NewCustomerListEntity(customers, c)
	controller.Response(c, 200, data)
}

// GetCustomer godoc
// @Tags Customer
// @Description meta value will be string or entities.FieldOption array
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param id path string true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.Customer}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /customer/:id [get]
func (ctl *CustomerController) GetCustomer(c *gin.Context) {

	customerID := c.Param("id")
	var err error

	v, ok := c.Get("FieldPermissions")
	if !ok {
		c.AbortWithStatusJSON(403, gin.H{"message": "Permission denied"})
		return
	}
	fieldPermissions := v.(string)

	levels := ""
	v, ok = c.Get("Levels")
	if ok {
		levels = v.(string)
	}

	customer, err := ctl.cusMgr.GetCustomer(c, customerID, fieldPermissions, levels)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get customer failed", err.Error())
		return
	}
	data := map[string]interface{}{}

	data["customer"] = entities.NewCustomerEntity(customer, c)
	controller.Response(c, 200, data)
}

// CreateCustomer godoc
// @Tags Customer
// @Description Add the string meta by "meta": {"meta1": "meta value"}, multiple meta by "meta": {"meta1": "optionId;optionId"}
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param payload body managers.CustomerCreateParam true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.Customer}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /customers [post]
func (ctl *CustomerController) Create(c *gin.Context) {

	v, ok := c.Get("AccountID")
	if !ok {
		c.AbortWithStatusJSON(403, gin.H{"message": "Permission denied"})
		return
	}
	accountID := v.(uuid.UUID)

	param := new(managers.CustomerCreateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	customer, err := ctl.cusMgr.Create(c, accountID, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "create customers failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["customer"] = entities.NewCustomerEntity(customer, c)
	controller.Response(c, 200, data)
}

// UpdateCustomer godoc
// @Tags Customer
// @Description Add the string meta by "meta": {"meta1": "meta value"}, multiple meta by "meta": {"meta1": "optionId;optionId"}
// @Accept json
// @Produce json
// @Param Authorization header string true " "
// @Param id path string true " "
// @Param payload body managers.CustomerUpdateParam true " "
// @Success 200 {object} swagger.APIResponse{data=swagger.Customer}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router /customer/:id [put]
func (ctl *CustomerController) Update(c *gin.Context) {

	v, ok := c.Get("AccountID")
	if !ok {
		c.AbortWithStatusJSON(403, gin.H{"message": "Permission denied"})
		return
	}
	accountID := v.(uuid.UUID)

	customerID := c.Param("id")
	param := new(managers.CustomerUpdateParam)
	err := controller.GetBody(c, param)
	if err != nil {
		return
	}
	customer, err := ctl.cusMgr.Update(c, accountID, customerID, param)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "create customers failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["customer"] = entities.NewCustomerEntity(customer, c)
	controller.Response(c, 200, data)
}
