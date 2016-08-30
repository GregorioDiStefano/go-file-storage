package log

import (
	"os"

	"github.com/Sirupsen/logrus"
)

type Entry logrus.Entry
type Fields logrus.Fields
type Level logrus.Level

const (
	// PanicLevel level, highest level of severity. Logs and then calls
	// panic with the message passed to Debug, Info, ...
	PanicLevel = logrus.PanicLevel
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even
	// if the logging level is set to Panic.
	FatalLevel = logrus.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be
	// noted.
	ErrorLevel = logrus.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = logrus.WarnLevel
	// InfoLevel level. General operational entries about what's going on
	// inside the application.
	InfoLevel = logrus.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose
	// logging.
	DebugLevel = logrus.DebugLevel
)

type Logger struct {
	// Inherit everything from logrus.Logger.
	logrus.Logger
}

// A global reference to our logger.
var Log *Logger

func init() {
	Log = NewLogger()
}

func NewLogger() *Logger {
	log := Logger{Logger: *logrus.New()}
	log.Out = os.Stderr
	return &log
}

// Mirror all the logrus API.

func ParseLevel(level string) (logrus.Level, error) {
	return logrus.ParseLevel(level)
}

func AddHook(hook logrus.Hook) {
	Log.Hooks.Add(hook)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func (self *Logger) Error(args ...interface{}) {
	self.error_impl(args...)
}

func (self *Logger) error_impl(args ...interface{}) {
	Log.Debug(args)
}

func (self *Logger) Info(args ...interface{}) {
	self.info_impl(args...)
}

func (self *Logger) info_impl(args ...interface{}) {
	self.Logger.Info(args)
}

func (self *Logger) Debug(args ...interface{}) {
	self.debug_impl(args...)
}

func (self *Logger) debug_impl(args ...interface{}) {
	self.Logger.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	Log.debugf_impl(format, args...)
}

func (self *Logger) Debugf(format string, args ...interface{}) {
	self.debugf_impl(format, args...)
}

func (self *Logger) debugf_impl(format string, args ...interface{}) {
	self.Logger.Debugf(format, args)
}

func Debugln(args ...interface{}) {
	Log.debugln_impl(args...)
}

func (self *Logger) Debugln(args ...interface{}) {
	self.debugln_impl(args...)
}

func (self *Logger) debugln_impl(args ...interface{}) {
	Log.Debugln(args)
}

func SetLevel(level logrus.Level) {
	Log.Level = level
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return Log.with_fields_impl(fields)
}

func (self *Logger) with_fields_impl(fields map[string]interface{}) *logrus.Entry {
	entry := self.WithFields(fields)
	return entry
}
