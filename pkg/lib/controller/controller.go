package controller

import (
	"encoding/json"
	"strings"

	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/gin-gonic/gin"
)

func GetQueryParam(c *gin.Context, v interface{}) error {

	return c.BindQuery(v)
}

func GetBody(c *gin.Context, v interface{}) error {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return err
	}
	if err := json.Unmarshal(jsonData, v); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return err
	}
	return nil
}

func HasPermission(c *gin.Context, permission string) bool {
	v, ok := c.Get("Permissions")
	if !ok {
		return false
	}
	userPermission := v.(string)
	if userPermission == _const.ROLE_SUPER {
		return true
	}
	// userPermission := encrypt.ASEDecrypt(v.(string), config.GetConfig().ASEKey)
	if !strings.Contains(userPermission, permission) {
		return false
	}
	return true
}

func Response(c *gin.Context, status int, data map[string]interface{}) {
	c.JSON(status, gin.H{
		"status": status,
		"data":   data,
	})
}

func ErrorResponse(c *gin.Context, status int, errorCode, message, detail string) {

	c.AbortWithStatusJSON(status, gin.H{
		"error":   errorCode,
		"message": message,
		"detail":  detail,
	})
}
