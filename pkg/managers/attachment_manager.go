package managers

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/service/firebase"
	"dmglab.com/mac-crm/pkg/util"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAttachmentManager interface {
	Upload(ctx context.Context, file *multipart.FileHeader, dir ...string) (*models.Attachment, error)
	Remove(ctx context.Context, id string) error
	GetAttachment(ctx context.Context, id string) (*models.Attachment, error)
	GetAttachments(ctx context.Context, dir ...string) ([]*models.Attachment, error)
}

type AttachmentManager struct {
	config *config.Config
}

func GetAttachmentManager() IAttachmentManager {
	return &AttachmentManager{
		config: config.GetConfig(),
	}
}

func (m *AttachmentManager) Upload(ctx context.Context, multipartFile *multipart.FileHeader, dir ...string) (*models.Attachment, error) {
	f, err := multipartFile.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var path string
	id := uuid.New()
	if len(dir) > 0 {
		path = filepath.Join(dir...)
		path = filepath.Join(m.config.FileStorage.LocalPath, m.config.CompanyID, path)
	} else {
		path = filepath.Join(m.config.FileStorage.LocalPath, m.config.CompanyID, "attachments")
	}
	if m.config.FileStorage.Driver == _const.LOCAL_STORAGE {
		_ = os.MkdirAll(path, os.ModePerm)
		path = filepath.Join(path, id.String())
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		_, err = io.Copy(file, f)
		if err != nil {
			return nil, err
		}
	} else if m.config.FileStorage.Driver == _const.GCP_STORAGE {
		err = firebase.GetFirebaseService(ctx).UploadFile(path, id.String(), f)
		if err != nil {
			return nil, err
		}
	}
	attachment := new(models.Attachment)
	attachment.ID = id
	attachment.FileName = multipartFile.Filename
	attachment.Size = int(multipartFile.Size)
	attachment.Path = path
	attachment.StorageDriver = m.config.FileStorage.Driver
	if ok := multipartFile.Header["Content-Type"]; ok != nil && len(multipartFile.Header["Content-Type"]) > 0 {
		attachment.MimeType = multipartFile.Header["Content-Type"][0]
	}

	return attachment, util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		return tx.Model(attachment).Create(attachment).Error
	})
}

func (m *AttachmentManager) Remove(ctx context.Context, id string) error {
	return util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		attachment := new(models.Attachment)
		attID, err := uuid.Parse(id)
		if err != nil {
			return err
		}
		attachment.ID = attID
		return tx.Model(attachment).Delete(attachment).Error
	})
}

func (m *AttachmentManager) GetAttachment(ctx context.Context, id string) (*models.Attachment, error) {
	attachment := new(models.Attachment)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		return tx.Model(attachment).First(attachment, "id = ?", id).Error
	})
	if err != nil {
		return nil, err
	}
	if m.config.FileStorage.Driver == _const.LOCAL_STORAGE {
		file, err := os.OpenFile(attachment.Path, os.O_RDONLY, os.ModePerm)
		if os.IsNotExist(err) {
			return nil, err
		}
		defer file.Close()
		attachment.Reader = file
		return attachment, err
	} else if m.config.FileStorage.Driver == _const.GCP_STORAGE {

		path := filepath.Join(attachment.Path, attachment.ID.String())
		file, err := firebase.GetFirebaseService(ctx).DownloadFile(path)
		if err != nil {
			return nil, err
		}
		attachment.Reader = file
		return attachment, err
	}
	return nil, nil
}

func (m *AttachmentManager) GetAttachments(ctx context.Context, dir ...string) ([]*models.Attachment, error) {

	var path string
	if len(dir) > 0 {
		path = filepath.Join(dir...)
		path = filepath.Join(m.config.FileStorage.LocalPath, m.config.CompanyID, path)
	} else {
		path = filepath.Join(m.config.FileStorage.LocalPath, m.config.CompanyID, "attachments")
	}
	attachments := []*models.Attachment{}
	return attachments, util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		return tx.Find(&attachments, `path LIKE ?`, path+"%").Error
	})
}
