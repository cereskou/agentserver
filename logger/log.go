package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var _log *logrus.Logger

func init() {
	_log = logrus.New()
	_log.Formatter = new(logrus.TextFormatter)
	_log.Level = logrus.FatalLevel
	_log.Out = os.Stdout
}

//GetLevel -
func GetLevel(log string) int {
	var lvl int = 0
	switch strings.ToLower(log) {
	case "debug":
		lvl = 1
	case "trace":
		lvl = 0
	case "info":
		lvl = 2
	case "warn":
		lvl = 3
	case "error":
		lvl = 4
	case "fatal":
		lvl = 5
	default:
		lvl = 4
	}
	return lvl
}

//SetLevel -
func SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		_log.Level = logrus.DebugLevel
	case "trace":
		_log.Level = logrus.TraceLevel
	case "info":
		_log.Level = logrus.InfoLevel
	case "warn":
		_log.Level = logrus.WarnLevel
	case "error":
		_log.Level = logrus.ErrorLevel
	case "fatal":
		_log.Level = logrus.FatalLevel
	default:
		_log.Level = logrus.FatalLevel
	}
}

//Println -
func Println(args ...interface{}) {
	fmt.Println(args...)
}

//Printf -
func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

//Error -
func Error(args ...interface{}) {
	_log.Error(args...)
}

//Errorf -
func Errorf(format string, args ...interface{}) {
	_log.Errorf(format, args...)
}

//Debug -
func Debug(args ...interface{}) {
	_log.Debug(args...)
}

//Debugf -
func Debugf(format string, args ...interface{}) {
	_log.Debugf(format, args...)
}

//Info -
func Info(args ...interface{}) {
	_log.Info(args...)
}

//Infof -
func Infof(format string, args ...interface{}) {
	_log.Infof(format, args...)
}

//Trace -
func Trace(args ...interface{}) {
	_log.Trace(args...)
}

//Tracef -
func Tracef(format string, args ...interface{}) {
	_log.Tracef(format, args...)
}

//Warn -
func Warn(args ...interface{}) {
	_log.Warn(args...)
}

//Warnf -
func Warnf(format string, args ...interface{}) {
	_log.Warnf(format, args...)
}

//Fatal -
func Fatal(args ...interface{}) {
	_log.Fatal(args...)
}

//Fatalf -
func Fatalf(format string, args ...interface{}) {
	_log.Fatalf(format, args...)
}
