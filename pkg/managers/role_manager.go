package managers

import (
	"context"

	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/util"
	"gorm.io/gorm"
)

type RoleCreateParam struct {
	Name             string `json:"name"`
	Permissions      string `json:"permissions" example:"ACCOUNT:U;CUSTOMER:U;ROLE:U"`
	FieldPermissions string `json:"field_permissions" example:"id_no;birth"`
}
type RoleUpdateParam struct {
	Name             string `json:"name"`
	Permissions      string `json:"permissions" example:"ACCOUNT:U;CUSTOMER:U;ROLE:U"`
	FieldPermissions string `json:"field_permissions" example:"id_no;birth"`
}
type IRoleManager interface {
	Create(ctx context.Context, param *RoleCreateParam) (*models.Role, error)
	Update(ctx context.Context, id string, param *RoleUpdateParam) (*models.Role, error)
	GetRoles(ctx context.Context) ([]*models.Role, error)
	GetRole(ctx context.Context, roleID string) (*models.Role, error)
}

type RoleManager struct {
	config *config.Config
}

func GetRoleManager() IRoleManager {
	return &RoleManager{
		config: config.GetConfig(),
	}
}

func (m *RoleManager) GetRoles(ctx context.Context) ([]*models.Role, error) {

	roles := []*models.Role{}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		var limit int
		l, ok := m.config.Setting[_const.ROW_LIMIT]
		if ok {
			limit = util.StrMustToInt(l)
		}
		err := tx.Model(&roles).Where("id != ? ", uuid.UUID{}).Limit(limit).Find(&roles).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (m *RoleManager) GetRole(ctx context.Context, roleID string) (*models.Role, error) {
	role := new(models.Role)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err := tx.First(role, "id = ?", roleID).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})
	if err != nil {
		//log
		return nil, err
	}
	return role, nil
}
func (m *RoleManager) Update(ctx context.Context, id string, param *RoleUpdateParam) (*models.Role, error) {
	roleNew := new(models.Role)

	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err := tx.First(roleNew, "id = ? AND permissions != ? ", id, _const.ROLE_SUPER).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		} else if err != nil {
			return err
		}
		if param != nil {
			roleNew.Name = param.Name
			roleNew.Permissions = param.Permissions
			roleNew.FieldPermissions = param.FieldPermissions
		}
		return tx.Save(roleNew).Error
	})
	if err != nil {
		return nil, err
	}
	return roleNew, nil
}
func (m *RoleManager) Create(ctx context.Context, param *RoleCreateParam) (*models.Role, error) {
	roleNew := new(models.Role)
	if param != nil {
		roleNew.Name = param.Name
		roleNew.Permissions = param.Permissions
		roleNew.FieldPermissions = param.FieldPermissions
	}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		return tx.Create(roleNew).Error

	})
	if err != nil {
		return nil, err
	}
	return roleNew, nil
}
