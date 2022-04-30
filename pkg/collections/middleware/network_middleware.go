package middleware

import (
	"dmglab.com/mac-crm/pkg/collections"
	"dmglab.com/mac-crm/pkg/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type INetworkMiddleware interface {
	CORSMiddleware() gin.HandlerFunc
	DBMiddleware() gin.HandlerFunc
}
type NetworkMiddleware struct {
	db *gorm.DB
}

func NewNetworkMiddleware() INetworkMiddleware {
	return &NetworkMiddleware{
		db: collections.GetCollection().DB,
	}
}
func (m *NetworkMiddleware) DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := m.db.Begin()
		defer func() {
			if r := recover(); r != nil {
				service.SysLog.Errorln(r)
				tx.Rollback()
			}
		}()
		tx.WithContext(c)
		c.Set("DB", tx)
		c.Next()
		if tx.Error != nil {
			tx.Rollback()
		} else if err := tx.Commit().Error; err != nil {
			service.SysLog.Error(err.Error())
		}
	}

}
func (m *NetworkMiddleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
