package firebase

import (
	"time"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type AuthEvent struct {
	Email    string `json:"email"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"metadata"`
	UID string `json:"uid"`
}

func (app *FirebaseApp) CreateAdminUser(tenantID, username, password, displayName string) (*auth.UserRecord, error) {
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(tenantID)
	if err != nil {
		return nil, err
	}

	user, err := tenantClient.CreateUser(app.Context, (&auth.UserToCreate{}).
		UID(uuid.New().String()).
		Email(username).
		EmailVerified(true).
		Password(password).
		DisplayName(displayName).
		Disabled(false))
	if err != nil {
		return nil, err
	}

	roleID := uuid.UUID{}
	customClaims := map[string]interface{}{
		"role_id": roleID.String(),
		"permission": map[string]interface{}{
			"SUPER": true,
		},
	}
	err = tenantClient.SetCustomUserClaims(app.Context, user.UID, customClaims)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (app *FirebaseApp) CreateTenant(tenant string) (*auth.Tenant, error) {

	return app.Auth.TenantManager.CreateTenant(app.Context, new(auth.TenantToCreate).
		AllowPasswordSignUp(true).
		DisplayName(tenant))

}
func (app *FirebaseApp) DeleteTenant(tenantID string) error {

	return app.Auth.TenantManager.DeleteTenant(app.Context, tenantID)

}
func (app *FirebaseApp) GetTenants() ([]*auth.Tenant, error) {
	tenants := []*auth.Tenant{}
	iter := app.Auth.TenantManager.Tenants(app.Context, "")
	if iter != nil {
		for {
			tenant, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			tenants = append(tenants, tenant)

		}
	}
	return tenants, nil
}

func (app *FirebaseApp) CreateUser(id, username, password, displayName string, claims map[string]interface{}) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		UID(id).
		Email(username).
		EmailVerified(true).
		// PhoneNumber("+15555550100").
		Password(password).
		DisplayName(displayName).
		// PhotoURL("http://www.example.com/12345678/photo.png").
		Disabled(false)
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(app.TenantID)
	if err != nil {
		return nil, err
	}
	user, err := tenantClient.CreateUser(app.Context, params)
	if err != nil {
		return nil, err
	}
	if claims == nil {
		return user, nil
	}
	return user, tenantClient.SetCustomUserClaims(app.Context, user.UID, claims)
}
func (app *FirebaseApp) UpdateUserProfile(id, displayName string, claims map[string]interface{}) (*auth.UserRecord, error) {
	params := (&auth.UserToUpdate{}).
		DisplayName(displayName)
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(app.TenantID)
	if err != nil {
		return nil, err
	}
	user, err := tenantClient.UpdateUser(app.Context, id, params)
	if err != nil {
		return nil, err
	}
	if claims == nil {
		return user, nil
	}
	return user, tenantClient.SetCustomUserClaims(app.Context, user.UID, claims)
}

func (app *FirebaseApp) DisableUser(id string) (*auth.UserRecord, error) {
	params := (&auth.UserToUpdate{}).
		Disabled(true)
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(app.TenantID)
	if err != nil {
		return nil, err
	}
	return tenantClient.UpdateUser(app.Context, id, params)
}

func (app *FirebaseApp) ActivateUser(id string) (*auth.UserRecord, error) {
	params := (&auth.UserToUpdate{}).
		Disabled(false)
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(app.TenantID)
	if err != nil {
		return nil, err
	}
	return tenantClient.UpdateUser(app.Context, id, params)
}

func (app *FirebaseApp) ChangePassword(id, password string) (*auth.UserRecord, error) {
	params := (&auth.UserToUpdate{}).
		Password("newPassword")
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(app.TenantID)
	if err != nil {
		return nil, err
	}
	return tenantClient.UpdateUser(app.Context, id, params)
}

func (app *FirebaseApp) GetUserByEmail(email string) (*auth.UserRecord, error) {

	return app.Auth.GetUserByEmail(app.Context, email)
}

func (app *FirebaseApp) GetUser(id string) (*auth.UserRecord, error) {
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(app.TenantID)
	if err != nil {
		return nil, err
	}
	return tenantClient.GetUser(app.Context, id)

}

func (app *FirebaseApp) GetUsers() ([]*auth.UserRecord, error) {
	users := []*auth.UserRecord{}
	var iter *auth.UserIterator
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(app.TenantID)
	if err != nil {
		return nil, err
	}
	iter = tenantClient.Users(app.Context, "")

	if iter != nil {
		for {
			user, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			users = append(users, &auth.UserRecord{
				UserInfo:               user.UserInfo,
				CustomClaims:           user.CustomClaims,
				Disabled:               user.Disabled,
				EmailVerified:          user.EmailVerified,
				ProviderUserInfo:       user.ProviderUserInfo,
				TokensValidAfterMillis: user.TokensValidAfterMillis,
				UserMetadata:           user.UserMetadata,
				TenantID:               user.TenantID,
			})

		}
	}
	return users, nil
}

func (app *FirebaseApp) SetPermissions(uid string, permissions map[string]interface{}) error {
	claims := map[string]interface{}{
		"permissions": permissions,
	}
	tenantClient, err := app.Auth.TenantManager.AuthForTenant(app.TenantID)
	if err != nil {
		return err
	}
	return tenantClient.SetCustomUserClaims(app.Context, uid, claims)
}

func (app *FirebaseApp) VerifyToken(token string) (*auth.Token, error) {
	// app.Auth.TenantManager.Tenant(app.Context,tenantID)
	authToken, err := app.Auth.VerifyIDToken(app.Context, token)
	if err != nil {
		return nil, err
	}
	return authToken, nil

}

func (app *FirebaseApp) CreateToken(uid string) (string, error) {
	return app.Auth.CustomTokenWithClaims(app.Context, uid, map[string]interface{}{})
}
