package middleware

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/dgrijalva/jwt-go"

	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/service"
	"github.com/gin-gonic/gin"
)

type IAuthMiddleware interface {
	AuthRequired() gin.HandlerFunc
	PermissionRequire(permission ...string) gin.HandlerFunc
}
type AuthMiddleware struct {
}

func NewAuthMiddleware() IAuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bearToken string
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer") {
			c.AbortWithStatusJSON(403, gin.H{"message": "invalid token"})
			return
		}
		if bearTokenStrings := strings.Split(authHeader, " "); len(bearTokenStrings) == 2 {
			bearToken = bearTokenStrings[1]
		} else {
			c.AbortWithStatusJSON(403, gin.H{"message": "invalid token"})
			return
		}

		acc := new(models.Account)
		// if v, ok := service.GetCache().Get(bearToken); ok {
		// 	acc = v.(*models.Account)
		// } else {
		token, err := jwt.Parse(bearToken, func(token *jwt.Token) (interface{}, error) {
			//Make sure that the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_TOKEN_SECRET")), nil
		})
		if err != nil {
			service.SysLog.Errorln(err.Error())
			c.AbortWithStatusJSON(403, gin.H{"message": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid && claims.VerifyExpiresAt(time.Now().Unix(), true) {
			v, ok := claims["tenant_id"].(string)
			if !ok {
				service.SysLog.Errorln("invalid Company ID")
				c.AbortWithStatusJSON(403, gin.H{"message": "invalid company id"})
				return
			}
			c.Set("TenantID", v)

			v, ok = claims["acc"].(string)
			if !ok {
				service.SysLog.Errorln("invalid Token")
				c.AbortWithStatusJSON(403, gin.H{"message": "invalid token"})
				return
			}

			err = json.Unmarshal([]byte(v), acc)
			if err != nil {
				service.SysLog.Errorln(err.Error())
				c.AbortWithStatusJSON(403, gin.H{"message": "invalid token"})
				return
			}
			// 	}
			// 	service.GetCache().Set(bearToken, acc, time.Minute*15)
		}
		c.Set("Account", acc)
		c.Set("FieldPermissions", acc.Role.FieldPermissions)
		c.Set("Levels", acc.Role.Levels)
		if acc.IsSystem {
			c.Set("Permissions", _const.ROLE_SUPER)
		} else {
			c.Set("Permissions", acc.Role.Permissions)
		}
		c.Next()

	}
}

func (m *AuthMiddleware) PermissionRequire(permission ...string) gin.HandlerFunc {

	return func(c *gin.Context) {
		v, ok := c.Get("Permissions")
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"message": "Permission denied"})
			return
		}
		userPermission := v.(string)
		// userPermission := encrypt.ASEDecrypt(v.(string), config.GetConfig().ASEKey)
		if userPermission == _const.ROLE_SUPER {
			c.Next()
		} else {
			for _, permission := range permission {
				if !strings.Contains(userPermission, permission) {
					c.AbortWithStatusJSON(403, gin.H{"message": "Permission denied"})
					return
				}
			}

			c.Next()
		}

	}
}
