package config

import (
	"crypto/rsa"
	"os"
	"strconv"

	"dmglab.com/mac-crm/pkg/service"
	"dmglab.com/mac-crm/pkg/util/encrypt"

	_const "dmglab.com/mac-crm/pkg/util/const"
)

type Config struct {
	CompanyID        string
	ASEKey           string
	DefaultListCount int
	DBDriver         string
	DB               *DBConfig
	App              *AppConfig
	FileStorage      *FileStorage
	Setting          map[string]string
}
type FileStorage struct {
	Driver    string
	LocalPath string
}

var conf *Config
var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
	Version  int
}

type AppConfig struct {
	Port    int
	Version string
}

func GetConfig() *Config {
	if conf != nil {
		return conf
	}
	licence := os.Getenv("LICENCE_KEY")
	if len(licence) == 0 || len(licence) < 40 {
		service.SysLog.Panicln("licence is invalid")
	}

	port := 8080
	if p := os.Getenv("PORT"); len(p) == 0 {
		service.SysLog.Errorln("the configuration for port is missing, initializing to port 8080")
	} else {
		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			service.SysLog.Errorln("failed initializing app port, use default port 8080")
			port = 8080
		}
	}
	storageDriver := os.Getenv("FILE_STORAGE")
	if len(storageDriver) == 0 {
		storageDriver = _const.LOCAL_STORAGE
	} else if storageDriver == _const.GCP_STORAGE {
		storageDriver = _const.GCP_STORAGE
	}

	key := encrypt.MD5Hash(licence[4:40])

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		service.SysLog.Errorln("failed initializing db port,db disable")
	}
	dbVersion, err := strconv.Atoi(os.Getenv("DB_VERSION"))
	if err != nil {
		service.SysLog.Errorln("failed initializing db version")
	}

	appConfig := &AppConfig{
		Port:    port,
		Version: os.Getenv("APP_VERSION"),
	}
	dbConfig := &DBConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		Database: os.Getenv("DB_NAME"),
		Version:  dbVersion,
	}
	fs := &FileStorage{
		Driver: storageDriver,
	}
	if storageDriver == _const.LOCAL_STORAGE {
		path := os.Getenv("LOCAL_STORAGE_PATH")
		if len(path) > 0 {
			fs.LocalPath = os.Getenv("LOCAL_STORAGE_PATH")
		} else {
			fs.LocalPath = "./asset"
			service.SysLog.Panicln("file storages to asset")
		}
	}

	conf = &Config{
		CompanyID:        os.Getenv("COMPANY_ID"),
		App:              appConfig,
		DefaultListCount: 5000,
		ASEKey:           key,
		DBDriver:         os.Getenv("DB_DRIVER"),
		DB:               dbConfig,
		FileStorage:      fs,
		Setting:          map[string]string{},
	}
	return conf

}
