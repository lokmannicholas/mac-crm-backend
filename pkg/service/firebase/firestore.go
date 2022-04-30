package firebase

import (
	"encoding/json"

	"dmglab.com/mac-crm/pkg/config"
)

type FireUserProfile struct {
	Permission map[string]bool `json:"permission"`
	TenantID   string          `json:"tenant_id"`
	RoleID     string          `json:"role_id"`
	UID        string          `json:"uid"`
}

func (app *FirebaseApp) GetUserFireProfile(email string) (*FireUserProfile, error) {
	snap, err := app.FireStore.Collection("users").
		Doc(email).Get(app.Context)
	if err != nil {
		return nil, err
	}

	m := snap.Data()["entries"]
	var entries []*FireUserProfile
	b, err := json.Marshal(&m)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &entries)
	if err != nil {
		return nil, err
	}
	var profile *FireUserProfile
	for _, p := range entries {
		if p.TenantID == config.GetConfig().CompanyID {
			profile = p
		}

	}
	return profile, nil
}
