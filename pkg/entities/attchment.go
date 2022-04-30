package entities

import (
	"dmglab.com/mac-crm/pkg/models"
)

type Attachment struct {
	ID        string `json:"id"`
	FileName  string `json:"file_name"`
	Size      int    `json:"size"`
	Path      string `json:"path"`
	Status    string `json:"status"`
	MimeType  string `json:"mime_type"`
	CreatedAt int64  `json:"created_at"`
}

func NewAttachmentEntity(attachment *models.Attachment) *Attachment {

	return &Attachment{
		ID:        attachment.ID.String(),
		FileName:  attachment.FileName,
		Size:      attachment.Size,
		Path:      attachment.Path,
		Status:    attachment.Status,
		MimeType:  attachment.MimeType,
		CreatedAt: attachment.CreatedAt.Unix(),
	}
}
func NewAttachmentListEntity(total int64, Attachments []*models.Attachment) *List {
	AttachmentList := make([]*Attachment, len(Attachments))
	for i, Attachment := range Attachments {
		AttachmentList[i] = NewAttachmentEntity(Attachment)
	}
	return &List{
		Total: total,
		Data:  AttachmentList,
	}

}
