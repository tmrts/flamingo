// Package flog implements logging utilities for flamingo
package flog

import "github.com/Sirupsen/logrus"

type Field func() (string, string)

func transform(fields []Field) logrus.Fields {
	logrusFields := logrus.Fields{}

	for _, f := range fields {
		k, v := f()

		logrusFields[k] = v
	}

	return logrusFields
}

func Debug(msg string, fields ...Field) {
	f := transform(fields)

	logrus.WithFields(f).Debug(msg)
}

func Info(msg string, fields ...Field) {
	f := transform(fields)

	logrus.WithFields(f).Info(msg)
}

func Warn(msg string, fields ...Field) {
	f := transform(fields)

	logrus.WithFields(f).Warning(msg)
}

func Error(msg string, fields ...Field) {
	f := transform(fields)

	logrus.WithFields(f).Error(msg)
}

func Fatal(msg string, fields ...Field) {
	f := transform(fields)

	logrus.WithFields(f).Fatal(msg)
}

func Panic(msg string, fields ...Field) {
	f := transform(fields)

	logrus.WithFields(f).Panic(msg)
}
