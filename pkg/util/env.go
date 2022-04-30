package util

import (
	"os"
	"time"

	"dmglab.com/mac-crm/pkg/service"
	"github.com/joho/godotenv"
)

func LoadENV() {
	var err error
	loc, err := time.LoadLocation(os.Getenv("TIME_ZONE"))
	if err != nil {
		time.Local = time.UTC
		service.SysLog.Errorln(err.Error())
	} else {
		time.Local = loc
	}

	runEnvironment := os.Getenv("ENV")

	service.SysLog.Println("running env file " + runEnvironment)
	if len(runEnvironment) > 0 {
		runEnvironment = ".env." + runEnvironment
		err := godotenv.Load(runEnvironment)
		if err != nil {
			service.SysLog.Errorln(err.Error())
		}
		service.SysLog.Errorln("running env file " + runEnvironment)
		return
	}
	err = godotenv.Load(".env")
	if err != nil {
		service.SysLog.Errorln(err.Error())
	}
	service.SysLog.Infoln("running env file .env")

}
