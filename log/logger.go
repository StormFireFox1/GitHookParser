package log

import (
	"github.com/sirupsen/logrus"
)

func Error(fields logrus.Fields, err error) {
	Logger.WithFields(fields).Error(err)
}

func Warn(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Warn(args)
}

func Info(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Info(args)
}

func Fatal(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Fatal(args)
}

func Panic(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Panic(args)
	panic(args)
}
