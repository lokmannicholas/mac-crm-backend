package service

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	_const "dmglab.com/mac-crm/pkg/util/const"

	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var auditLog *AuditLog

type AuditLog struct {
	FileName string
	*logrus.Logger
}

func GetAuditLogger() *AuditLog {
	date := time.Now().Format("2006_01_02")
	fileName := fmt.Sprintf("mac_pro_%s.log", date)

	if auditLog == nil || auditLog.Logger == nil {
		// log = logrus.New()
		customFormatter := new(logrus.TextFormatter)
		customFormatter.TimestampFormat = "2006-01-02 15:04:05"
		customFormatter.FullTimestamp = true
		// log.SetFormatter(customFormatter)
		auditLog = &AuditLog{
			FileName: "",
			Logger: &logrus.Logger{
				Out:       os.Stdout,
				Formatter: customFormatter,
			},
		}
	}
	path := os.Getenv("LOG_PATH")
	if len(path) > 0 {
		//file by date
		readFile := false
		if auditLog.FileName != fileName {
			auditLog.FileName = fileName
			readFile = true
		}

		if readFile {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic(err)
			}

			file, err := os.OpenFile(
				filepath.Join(path, fileName),
				os.O_CREATE|os.O_WRONLY|os.O_APPEND,
				0666)
			if err == nil {
				auditLog.Out = file
			} else {
				auditLog.Info("Failed to log to file, using default stderr")
			}
		}
	} else {
		auditLog.SetOutput(os.Stdout)
	}
	return auditLog
}

func (aud *AuditLog) InfoLog(ctx context.Context, action _const.AuditAction, arg ...interface{}) {

	auth := ctx.Value("Auth")
	if auth == nil {
		aud.Logger.Errorln("invalid account in context")
		aud.Logger.Logf(logrus.InfoLevel, "unauthorized action %s with content %+v", action, arg)
	} else {
		if b, err := json.Marshal(arg); err != nil {
			aud.Logger.Error(err)
			aud.Logger.Logf(logrus.InfoLevel, "unauthorized action %s with content %+v", action, arg)
		} else {
			acc := auth.(map[string]string)
			aud.Logf(logrus.InfoLevel, "%s(%s) excute action %s with content %s\n", acc["username"], acc["id"], action, string(b))
		}
	}
}

func (aud *AuditLog) SetOut(w io.Writer) *AuditLog {
	aud.Logger.Out = w
	return aud
}
