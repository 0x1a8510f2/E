package log

import (
	"os"
	"strings"

	logLib "github.com/sirupsen/logrus"
)

func Init(output *os.File) {
	logLib.SetOutput(output)
}

func SetLevel(levelString string) {
	levelString = strings.ToLower(levelString)
	switch levelString {
	case "trace":
		logLib.SetLevel(logLib.TraceLevel)
	case "debug":
		logLib.SetLevel(logLib.DebugLevel)
	case "info":
		logLib.SetLevel(logLib.InfoLevel)
	case "warn":
		logLib.SetLevel(logLib.WarnLevel)
	case "error":
		logLib.SetLevel(logLib.ErrorLevel)
	case "fatal":
		logLib.SetLevel(logLib.FatalLevel)
	case "panic":
		logLib.SetLevel(logLib.PanicLevel)
	default:
		logLib.SetLevel(logLib.InfoLevel)
	}
}

func Tracef(str string, args ...interface{}) {
	logLib.Tracef(str, args...)
}

func Debugf(str string, args ...interface{}) {
	logLib.Debugf(str, args...)
}

func Infof(str string, args ...interface{}) {
	logLib.Infof(str, args...)
}

func Warnf(str string, args ...interface{}) {
	logLib.Warnf(str, args...)
}

func Errorf(str string, args ...interface{}) {
	logLib.Errorf(str, args...)
}

func Fatalf(str string, args ...interface{}) {
	logLib.Fatalf(str, args...)
}

func Panicf(str string, args ...interface{}) {
	logLib.Panicf(str, args...)
}
