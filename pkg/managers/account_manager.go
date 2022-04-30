package managers

import (
	"context"

	_const "dmglab.com/mac-crm/pkg/util/const"

	"dmglab.com/mac-crm/pkg/util/encrypt"

	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"dmglab.com/mac-crm/pkg/service"
	"dmglab.com/mac-crm/pkg/util"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"gorm.io/gorm"
)

type AccountCreateParam struct {
	ID          *string ` json:"-"`
	DisplayName string  `validate:"required" json:"display_name"`
	Username    string  `validate:"required" json:"username"`
	Password    string  `validate:"required" json:"password"`
	RoleID      string  `validate:"required" json:"role_id"`
}
type AccountUpdateParam struct {
	DisplayName string `validate:"required" json:"display_name"`
	Username    string `validate:"required" json:"username"`
	Password    string `validate:"required" json:"password"`
	RoleID      string `validate:"required" json:"role_id"`
	Status      string `json:"status"`
}
type ChangePasswordParam struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}

type LoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type IAccountManager interface {
	Create(ctx context.Context, param *AccountCreateParam) (*models.Account, error)
	Update(ctx context.Context, id string, param *AccountUpdateParam) (*models.Account, error)
	Activate(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error
	CreateToken(ctx context.Context) (string, string, error)
	GetAccounts(ctx context.Context, ids ...string) ([]*models.Account, error)
	GetSchedulerAccount() *models.Account
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	Login(ctx context.Context, param *LoginParam) (*models.Account, error)
	Logout(ctx context.Context) error
	ChangePassword(ctx context.Context, id string, param *ChangePasswordParam) error
}

type AccountManager struct {
	config *config.Config
}

func GetAccountManager() IAccountManager {
	return &AccountManager{
		config: config.GetConfig(),
	}
}

/*Create a new account*/
func (m *AccountManager) Create(ctx context.Context, param *AccountCreateParam) (*models.Account, error) {
	var acc *models.Account
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		roleUUID, err := uuid.Parse(param.RoleID)
		if err != nil {
			return err
		}
		if param.ID != nil {
			id, err := uuid.Parse(*param.ID)
			if err != nil {
				return err
			}
			acc = &models.Account{ID: id, DisplayName: param.DisplayName, Username: param.Username, Password: param.Password, RoleID: roleUUID, IsSystem: false}
		} else {
			acc = &models.Account{DisplayName: param.DisplayName, Username: param.Username, Password: param.Password, RoleID: roleUUID, IsSystem: false}
		}

		return tx.Create(acc).Error
	})
	if err != nil {
		return nil, err
	}
	return acc, nil
}

/*update account*/
func (m *AccountManager) Update(ctx context.Context, id string, param *AccountUpdateParam) (*models.Account, error) {
	acc := new(models.Account)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		roleUUID, err := uuid.Parse(param.RoleID)
		if err != nil {
			return err
		}

		err = tx.Where("is_system = ?", false).First(acc, "id = ?", id).Error
		if err != nil {
			return err
		}
		acc.DisplayName = param.DisplayName
		acc.Username = param.Username
		acc.SetPassword(param.Password)
		acc.RoleID = roleUUID
		acc.Status = param.Status
		return tx.Save(acc).Error
	})
	if err != nil {
		return nil, err
	}
	return acc, nil
}

/*Disable an account*/
func (m *AccountManager) Disable(ctx context.Context, id string) error {
	var acc *models.Account
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		acc = new(models.Account)
		err := tx.Where("is_system = ?", false).First(acc, "id = ?", id).Error
		if err != nil {
			//log
			return err
		}
		acc.SetDisable()
		return tx.Save(acc).Error
	})
	if err != nil {
		//log
		return err
	}
	return nil
}
func (m *AccountManager) Activate(ctx context.Context, id string) error {
	acc := new(models.Account)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {

		err := tx.Where("is_system = ?", false).First(acc, "id = ? ", id).Error
		if err != nil {
			return err
		}
		acc.SetActive()
		return tx.Save(acc).Error
	})
	if err != nil {
		//log
		return err
	}
	return nil
}

func (m *AccountManager) GetAccounts(ctx context.Context, ids ...string) ([]*models.Account, error) {
	accounts := []*models.Account{}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		tx = tx.Preload("Role")
		if len(ids) > 0 {
			tx = tx.Where("id IN (?)", ids)
		}
		err := tx.Find(&accounts, "is_system = ?", false).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		return err

	})

	return accounts, err
}

func (m *AccountManager) GetSchedulerAccount() *models.Account {
	return &models.Account{
		ID:          uuid.UUID{},
		DisplayName: _const.ACC_SCHEDULER,
		Username:    "scheduler",
		Status:      "Active",
		IsSystem:    true,
	}
}

func (m *AccountManager) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	acc := new(models.Account)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err := tx.First(acc, "id = ?", acc.ID).Error
		if err != nil {
			return err
		}

		return tx.First(&acc.Role, "id = ?", acc.RoleID).Error
	})
	if err != nil {
		//log
		return nil, err
	}
	return acc, nil
}

/*Login return base64 encode*/
func (m *AccountManager) CreateToken(ctx context.Context) (string, string, error) {
	var err error
	acc := new(models.Account)
	if acc, err = util.GetCtxAccount(ctx); err != nil {
		return "", "", err
	}
	err = util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err = tx.First(acc, "id = ?", acc.ID).Error
		if err != nil {
			return err
		}

		return tx.First(&acc.Role, "id = ?", acc.RoleID).Error
	})
	if err != nil {
		return "", "", err
	}

	atClaims := jwt.MapClaims{
		"tenant_id":  config.GetConfig().CompanyID,
		"authorized": true,
		"acc":        acc.ToJSON(),
		"exp":        time.Now().Add(time.Hour * 10).Unix(),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET")))
	if err != nil {
		return "", "", err
	}
	rtClaims := jwt.MapClaims{
		"expire": time.Now().Add(time.Hour * 1).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func (m *AccountManager) Login(ctx context.Context, param *LoginParam) (*models.Account, error) {
	acc := new(models.Account)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		if err := tx.Preload("Role").First(acc, "is_system = ? AND username = ? AND password = ?", false, param.Username, acc.SetPassword(param.Password)).Error; err != nil {
			return err
		}

		if err := acc.LastLogin.Scan(time.Now()); err != nil {
			return err
		}
		if err := tx.Save(acc).Error; err != nil {
			return err
		}
		//audit log
		service.GetAuditLogger().Logf(logrus.InfoLevel, "%s(%s) login\n", acc.Username, acc.ID.String())
		return nil
	})

	return acc, err
}

func (m *AccountManager) Logout(ctx context.Context) error {

	return util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		acc, err := util.GetCtxAccount(ctx)
		if err != nil {
			return err
		}
		service.GetAuditLogger().Logf(logrus.InfoLevel, "%s(%s) logout\n", acc.Username, acc.ID.String())
		return nil
	})

}

func (m *AccountManager) ChangePassword(ctx context.Context, id string, param *ChangePasswordParam) error {

	return util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		newPassword := encrypt.MD5Hash(param.Password)
		oldPassword := encrypt.MD5Hash(param.OldPassword)
		return tx.Model(&models.Account{}).Where("id = ? AND password = ?", id, oldPassword).Update("password", newPassword).Error
	})

}
