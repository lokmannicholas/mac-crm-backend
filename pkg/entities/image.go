package entities

import "dmglab.com/mac-crm/pkg/models"

type Image struct {
	ID        string `json:"id"`
	ImageType string `json:"image_type"`
	Data      string `json:"data"`
	RefID     string `json:"ref_id"`
}

func NewImageEntity(img *models.Image) *Image {

	r := &Image{
		ID:        img.ID.String(),
		ImageType: img.ImageType,
		Data:      string(img.Data),
		RefID:     img.RefID.String(),
	}
	return r
}

func NewImageListEntity(images []*models.Image) []*Image {
	imageList := make([]*Image, len(images))
	for i, image := range images {
		imageList[i] = NewImageEntity(image)
	}
	return imageList

}
