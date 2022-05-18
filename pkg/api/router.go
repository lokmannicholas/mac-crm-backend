package api

import (
	"errors"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	docs "dmglab.com/mac-crm/docs"
	DBMiddleware "dmglab.com/mac-crm/pkg/collections/middleware"
	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/lib/auth/middleware"
	_const "dmglab.com/mac-crm/pkg/util/const"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter() *gin.Engine {
	//middleWare
	if len(config.GetConfig().CompanyID) == 0 {
		panic(errors.New("invalid company id"))
	}
	netMidWare := DBMiddleware.NewNetworkMiddleware()
	router := gin.Default()
	router.Use(netMidWare.DBMiddleware())
	router.Use(netMidWare.CORSMiddleware())

	//swagger
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/" + config.GetConfig().CompanyID
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r := router.Group("/api/" + config.GetConfig().CompanyID)
	{

		nonAuthPath(r)
		apiPath(r)

	}

	return router

}
func nonAuthPath(r *gin.RouterGroup) {
	authMethod := config.GetConfig().Auth
	if authMethod == "onpremise" {
		accCtl := NewAccountController()
		r.POST("/auth/login", accCtl.Login)
		r.POST("/auth/logout", accCtl.Logout)
	} else if authMethod == "firebase" {
	}
	r.GET("/config", gin.HandlerFunc(func(c *gin.Context) {
		conf := config.GetConfig()
		licenceValid := "APPROVED"
		licence := os.Getenv("LICENCE_KEY")
		if len(licence) == 0 || len(licence) < 40 {
			licenceValid = "REJECTED"
		}
		c.JSON(200, gin.H{
			"status": 200,
			"data": map[string]string{
				"db":           conf.DB.Host,
				"db-version":   strconv.Itoa(conf.DB.Version),
				"db-driver":    conf.DBDriver,
				"storage":      conf.FileStorage.Driver,
				"storage-path": conf.FileStorage.LocalPath,
				"app-version":  conf.App.Version,
				"licence":      licenceValid,
				"time":         time.Now().Format("2006-01-02T15:04:05Z07:00"),
			},
		})
	}))
}
func apiPath(r *gin.RouterGroup) {
	accMidWare := middleware.NewAuthMiddleware()
	r.Use(accMidWare.AuthRequired())
	authMethod := config.GetConfig().Auth
	if authMethod == "onpremise" {
		NewAccountController().SetRouter(r)
	} else if authMethod == "firebase" {

	}
	NewSettingController().SetRouter(r)
	NewCustomFieldController().SetRouter(r)
	NewFieldOptionController().SetRouter(r)
	attCtl := NewAttachmentController()
	rolCtl := NewRoleController()
	cusCtl := NewCustomerController()
	pmsCtl := NewPermissionController()

	// socket := ws.NewWebSocketController()
	//websocket

	// r.GET("/notification", socket.Notification)
	//restful
	r.GET("/permissions", pmsCtl.GetPermissions)

	r.GET("/roles", accMidWare.PermissionRequire(_const.PERMISSION_ROLE.Read()), rolCtl.GetRoles)
	r.POST("/roles", accMidWare.PermissionRequire(_const.PERMISSION_ROLE.Create()), rolCtl.Create)
	r.PUT("/role/:id", accMidWare.PermissionRequire(_const.PERMISSION_ROLE.Update()), rolCtl.Update)

	r.GET("/attachment/:id", attCtl.GetAttachment)

	r.GET("/customers", accMidWare.PermissionRequire(_const.PERMISSION_CUSTOMER.Read()), cusCtl.GetCustomers)
	r.POST("/customers", accMidWare.PermissionRequire(_const.PERMISSION_CUSTOMER.Create()), cusCtl.Create)

	customer := r.Group("/customer", accMidWare.PermissionRequire(_const.PERMISSION_CUSTOMER.Read()))
	{
		customer.GET("/:id", accMidWare.PermissionRequire(_const.PERMISSION_CUSTOMER.Read()), cusCtl.GetCustomer)
		customer.PUT("/:id", accMidWare.PermissionRequire(_const.PERMISSION_CUSTOMER.Update()), cusCtl.Update)
	}

}

func ReverseProxy(target string) gin.HandlerFunc {
	url, _ := url.Parse(target)
	// (err)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
