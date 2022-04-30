package collections

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/service"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"os"
)

type DatabaseConnection struct {
	*gorm.DB
}

var dbc *DatabaseConnection

func GetCollection() *DatabaseConnection {
	if dbc == nil {
		dbc = NewCollection()
	}
	return dbc
}

func (db *DatabaseConnection) Close() {
	if db != nil {
		db.Close()
	}
}
func NewCollection() *DatabaseConnection {

	conf := config.GetConfig()
	var dia gorm.Dialector
	prepareStmt := false
	switch conf.DBDriver {
	case "CLOUDMYSQL":
		{
			dsn := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
				conf.DB.Username, conf.DB.Password, os.Getenv("INSTANCE_CONNECTION_NAME"), conf.DB.Database,
				url.QueryEscape(os.Getenv("TIME_ZONE")))

			dia = mysql.Open(dsn)
			prepareStmt = false
			break
		}
	case "MSSQL":
		{
			dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
				conf.DB.Username,
				conf.DB.Password,
				conf.DB.Host,
				conf.DB.Port,
				conf.DB.Database,
			)
			dia = sqlserver.Open(dsn)
			prepareStmt = true
			break
		}
	case "MYSQL":
		{
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
				conf.DB.Username,
				conf.DB.Password,
				conf.DB.Host,
				conf.DB.Port,
				conf.DB.Database,
				url.QueryEscape(os.Getenv("TIME_ZONE")),
			)

			dia = mysql.Open(dsn)
			prepareStmt = false
			break
		}
	default:
		{
			dia = sqlite.Open("gorm.db")
		}
	}
	db, err := gorm.Open(dia, &gorm.Config{
		Logger:      newLogger(),
		PrepareStmt: prepareStmt,
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return &DatabaseConnection{db}

}

func newLogger() logger.Interface {
	l := logger.New(
		service.SysLog,
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	return l.LogMode(logger.Info)
}

type AppSetting struct {
	AppVersion string
	LicenceKey string
	ASEKey     string
	DBDriver   string
	DBVersion  int
}

func (d *DatabaseConnection) getAppSetting(migration *gorm.DB) (*AppSetting, error) {

	if !migration.Migrator().HasTable(&models.App{}) {
		migration.AutoMigrate(&models.App{})
		settings := []models.App{
			{
				Setting: "APP_VERSION",
				Value:   []byte(os.Getenv("APP_VERSION")),
				Version: 0,
			},
			{
				Setting: "DB_VERSION",
				Value:   []byte(os.Getenv("DB_VERSION")),
				Version: 0,
			},
			{
				Setting: "ASE_KEY",
				Value:   []byte(config.GetConfig().ASEKey),
				Version: 0,
			},
			{
				Setting: "LICENCE_KEY",
				Value:   []byte(os.Getenv("LICENCE_KEY")),
				Version: 0,
			},
			{
				Setting: "DB_DRIVER",
				Value:   []byte(os.Getenv("DB_DRIVER")),
				Version: 0,
			},
		}

		if err := migration.Transaction(func(tx *gorm.DB) error {
			return tx.CreateInBatches(&settings, len(settings)).Error
		}); err != nil {
			return nil, err
		}
		dbVersion, err := strconv.Atoi(os.Getenv("DB_VERSION"))
		if err != nil {
			return nil, err
		}
		return &AppSetting{
			AppVersion: os.Getenv("APP_VERSION"),
			LicenceKey: os.Getenv("LICENCE_KEY"),
			ASEKey:     config.GetConfig().ASEKey,
			DBDriver:   os.Getenv("DB_DRIVER"),
			DBVersion:  dbVersion,
		}, nil
	} else if migration.Migrator().HasTable(&models.App{}) {
		results := map[string]string{}
		err := migration.Transaction(func(tx *gorm.DB) error {
			settings := []*models.App{}
			err := tx.Find(&settings).Error
			if err != nil {
				return err
			}
			for _, setting := range settings {
				results[setting.Setting] = string(setting.Value)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if len(results) > 0 {
			dbVersion, err := strconv.Atoi(results["DB_VERSION"])
			if err != nil {
				return nil, err
			}
			return &AppSetting{
				AppVersion: results["APP_VERSION"],
				LicenceKey: results["LICENCE_KEY"],
				ASEKey:     results["ASE_KEY"],
				DBDriver:   results["DB_DRIVER"],
				DBVersion:  dbVersion,
			}, err
		}
	}
	return nil, errors.New("app setting errors")
}
func (d *DatabaseConnection) Migrate(f func(migration *gorm.DB) error) {
	if isMigra, err := strconv.ParseBool(os.Getenv("DB_MIGRATION")); err != nil {
		return
	} else if !isMigra {
		return
	}

	migration := d.Session(&gorm.Session{
		PrepareStmt: false,
		Logger:      logger.Default.LogMode(logger.Info),
		SkipHooks:   true})
	appSetting, err := d.getAppSetting(migration)
	if err != nil {
		service.SysLog.Panicln(err.Error())
	}
	service.SysLog.Infof("current database version %d\n", appSetting.DBVersion)
	if config.GetConfig().DB.Version != 0 && appSetting.DBVersion >= config.GetConfig().DB.Version {
		return
	}
	if err := f(migration); err == nil {
		if version := config.GetConfig().DB.Version; version > 0 {
			if err := migration.Model(&models.App{}).Where("`setting` = ?", "DB_VERSION").Update("value", strconv.Itoa(version)).Error; err != nil {
				service.SysLog.Fatalln(err)
			}
		}
	}

}
