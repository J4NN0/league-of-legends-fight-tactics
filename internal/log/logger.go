package log

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	name string
	l    *logrus.Logger
}

func New(name string) *Logger {
	return &Logger{
		name: name,
		l:    logrus.New(),
	}
}

func (l *Logger) Printf(fmt string, args ...interface{}) {
	l.l.WithFields(logrus.Fields{"name": l.name}).Printf(fmt, args...)
}

func (l *Logger) Warningf(fmt string, args ...interface{}) {
	l.l.WithFields(logrus.Fields{"name": l.name}).Warningf(fmt, args...)
}

func (l *Logger) Fatalf(fmt string, args ...interface{}) {
	l.l.WithFields(logrus.Fields{"name": l.name}).Fatalf(fmt, args...)
}
