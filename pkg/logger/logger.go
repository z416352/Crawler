package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func new(fieldsOrder []string) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	log.SetReportCaller(true)
	log.Formatter = &nested.Formatter{
		FieldsOrder:     []string{"component", "category"},
		NoFieldsSpace:   true,
		HideKeys:        true,
		TimestampFormat: time.RFC850,
		CallerFirst:     true,

		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf(" [%s:%d][%s()]", path.Base(f.File), f.Line, funcName)
		},
	}
	return log
}

var (
	Log        *logrus.Logger
	MainLog    *logrus.Entry
	CrawlerLog *logrus.Entry
	UtilsLog   *logrus.Entry
	DBLog      *logrus.Entry
)

func init() {
	fieldsOrder := []string{
		"package",
		"category",
	}

	Log = new(fieldsOrder)
	MainLog = Log.WithField("package", "main")
	CrawlerLog = Log.WithField("package", "crawler")
	UtilsLog = Log.WithField("package", "utils")
	DBLog = Log.WithField("package", "DB")
}
