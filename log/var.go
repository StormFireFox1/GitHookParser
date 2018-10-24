package log

import (
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Logger is the logger for Git-Hook-Parser.
//
// It uses Lumberjack to roll logfiles and logrus as the backend. All logs can be parsed with Logstash.
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	Logger.Out = &lumberjack.Logger{
		Filename:   "log/hook-parser.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     365, // days
		Compress:   true,
	}

	Logger.Formatter = &logrus.JSONFormatter{}
	Logger.SetLevel(logrus.InfoLevel)
}
