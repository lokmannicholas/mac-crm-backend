package api

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/entities"
	controller "dmglab.com/mac-crm/pkg/lib/controller"
	"dmglab.com/mac-crm/pkg/managers"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

type IAttachmentController interface {
	GetAttachment(c *gin.Context)
	GetAttachments(c *gin.Context)
	Upload(c *gin.Context)
	Delete(c *gin.Context)
}
type AttachmentController struct {
	attMgr    managers.IAttachmentManager
	validator *validator.Validate
}

func NewAttachmentController() IAttachmentController {
	return &AttachmentController{
		attMgr:    managers.GetAttachmentManager(),
		validator: validator.New(),
	}
}

// GetAttachment godoc
// @Tags Attachment
// @Accept json
// @Produce json
// @param 		Authorization header string true "Authorization"
// @Success      200
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure      500  {object}  swagger.APIInternalServerError
// @Router       /attachment/:id [get]
func (ctl *AttachmentController) GetAttachment(c *gin.Context) {
	attID := c.Param("id")

	c.Stream(func(w io.Writer) bool {
		attachment, err := ctl.attMgr.GetAttachment(c, attID)
		if err != nil {
			controller.ErrorResponse(c, 500, "000000", "get Attachments failed", err.Error())
			return false
		}
		c.Header("Content-Type", attachment.MimeType)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", attachment.FileName))
		_, err = io.Copy(w, attachment.Reader)
		if err != nil {
			controller.ErrorResponse(c, 500, "000000", "get Attachments failed", err.Error())
			return false
		}
		return false
	})
}

func (ctl *AttachmentController) Delete(c *gin.Context) {
	attID := c.Param("attachmentID")
	err := ctl.attMgr.Remove(c, attID)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get Attachments failed", err.Error())
		return
	}
	controller.Response(c, 200, map[string]interface{}{
		"status": "OK",
	})

}

// Upload
// @Summary      Upload attachment
// @Description  Upload attachment
// @Tags         Customer
// @Accept multipart/form-data
// @Produce      json
// @param Authorization header string true "Authorization"
// @Param file formData file true "file"
// @Success 200 {object} swagger.APIResponse{data=entities.Attachment}
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router       /customer/:id/attachments [post]
func (ctl *AttachmentController) Upload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get retrieve reference file", err.Error())
		return
	}
	path := strings.TrimPrefix(c.Request.URL.Path, filepath.Join("/api", config.GetConfig().CompanyID))
	paths := strings.Split(path, "/")
	att, err := ctl.attMgr.Upload(c, file, paths...)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "upload attachment failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["attachment"] = entities.NewAttachmentEntity(att)
	controller.Response(c, 200, data)
}

// GetAttachments godoc
// @Tags Customer
// @Accept json
// @Produce json
// @param 		Authorization header string true "Authorization"
// @Success      200
// @Failure 403 {object} swagger.APIForbiddenError
// @Failure 500 {object} swagger.APIInternalServerError
// @Router       /customer/:id/attachments [get]
func (ctl *AttachmentController) GetAttachments(c *gin.Context) {

	path := strings.TrimPrefix(c.Request.URL.Path, filepath.Join("/api", config.GetConfig().CompanyID))
	paths := strings.Split(path, "/")
	attachments, err := ctl.attMgr.GetAttachments(c, paths...)
	if err != nil {
		controller.ErrorResponse(c, 500, "000000", "get Attachments failed", err.Error())
		return
	}
	data := map[string]interface{}{}
	data["attachments"] = entities.NewAttachmentListEntity(int64(len(attachments)), attachments)
	controller.Response(c, 200, data)
}
