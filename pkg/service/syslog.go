package service

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var SysLog = &logrus.Logger{
	Out:          os.Stdout,
	Level:        logrus.InfoLevel,
	Formatter:    &CustomFormate{},
	ReportCaller: true,
}

type CustomFormate struct{}

var levelList = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

func (mf *CustomFormate) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	level := levelList[int(entry.Level)]
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1]
	timeForm := entry.Time.Format("2006-01-02 15:04:05,678")
	switch entry.Level {
	case logrus.PanicLevel:
		{
			b.WriteString(fmt.Sprintf("%s - [%s] - %s:%d %s\n",
				timeForm, level, fileName,
				entry.Caller.Line, entry.Message))

		}
	case logrus.FatalLevel:
		{
			b.WriteString(fmt.Sprintf("%s - [%s] - %s:%d %s\n",
				timeForm, level, fileName,
				entry.Caller.Line, entry.Message))

		}
	case logrus.ErrorLevel:
		{
			b.WriteString(fmt.Sprintf("%s - [%s] - %s:%d %s\n",
				timeForm, level, fileName,
				entry.Caller.Line, entry.Message))

		}
	case logrus.DebugLevel:
		{
			b.WriteString(fmt.Sprintf("%s - [%s] - %s:%d %s\n",
				timeForm, level, fileName,
				entry.Caller.Line, entry.Message))

		}
	case logrus.TraceLevel:
		{
			b.WriteString(fmt.Sprintf("%s - [%s] - %s:%d %s\n",
				timeForm, level, fileName,
				entry.Caller.Line, entry.Message))

		}
	default:
		{
			b.WriteString(fmt.Sprintf("%s - [%s] - %s\n",
				timeForm, level, entry.Message))
		}
	}
	return b.Bytes(), nil
}
