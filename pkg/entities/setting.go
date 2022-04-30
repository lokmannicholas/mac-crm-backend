package entities

import (
	"dmglab.com/mac-crm/pkg/models"
)

type Setting struct {
	Setting string `json:"setting"`
	Value   string `json:"value"`
	Version int    `json:"version"`
}

func NewSettingEntity(app *models.App) *Setting {
	if app == nil {
		return &Setting{}
	}
	c := &Setting{
		Setting: app.Setting,
		Value:   string(app.Value),
		Version: app.Version,
	}

	return c
}

func NewSettingListEntity(total int64, apps []*models.App) *List {
	appsList := make([]*Setting, len(apps))
	for i, app := range apps {
		appsList[i] = NewSettingEntity(app)
	}
	return &List{
		Total: total,
		Data:  appsList,
	}

}
