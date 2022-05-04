package util

import (
	"context"

	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/service"

	"errors"

	"gorm.io/gorm"
)

func SetCtxDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, "DB", db)
}

func GetCtxDB(ctx context.Context) (*gorm.DB, error) {
	db := ctx.Value("DB").(*gorm.DB)
	if db != nil {
		return db, nil
	}
	return nil, errors.New("no database in context")
}

func GetCtxTx(ctx context.Context, f func(tx *gorm.DB) error) error {
	db := ctx.Value("DB").(*gorm.DB)
	if db == nil {
		return errors.New("no database in context")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		err := f(tx)
		if err != nil {
			service.SysLog.Errorln(err)
			if err != gorm.ErrRecordNotFound {
				tx.Rollback()
			}
			return err
		}
		return err
	})
}

func SetCtxAccount(ctx context.Context, acc *models.Account) context.Context {
	return context.WithValue(ctx, "Account", acc)
}
func GetCtxAccount(ctx context.Context) (*models.Account, error) {
	acc := ctx.Value("Account").(*models.Account)
	if acc != nil {
		return acc, nil
	}
	return nil, errors.New("no account in context")
}
