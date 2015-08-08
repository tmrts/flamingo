// Package flog implements logging utilities for flamingo
package flog

import "github.com/Sirupsen/logrus"

type Parameter interface {
	Convert() map[string]interface{}
}

type Fields struct {
	Event string
	Error error
}

func (f Fields) Convert() map[string]interface{} {
	fields := map[string]interface{}{}

	if f.Event != "" {
		fields["event"] = f.Event
	}

	if f.Error != nil {
		fields["error"] = f.Error
	}

	return fields
}

type Details map[string]interface{}

func (d Details) Convert() map[string]interface{} {
	return d
}

type DebugFields map[string]interface{}

func (d DebugFields) Convert() map[string]interface{} {
	// TODO(tmrts): Handle debug information
	return map[string]interface{}{}
}

func transform(params []Parameter) logrus.Fields {
	logrusFields := logrus.Fields{}

	for _, p := range params {
		fieldMap := p.Convert()

		for k, v := range fieldMap {
			logrusFields[k] = v
		}
	}

	return logrusFields
}

func Debug(msg string, params ...Parameter) {
	f := transform(params)

	logrus.WithFields(f).Debug(msg)
}

func Info(msg string, params ...Parameter) {
	f := transform(params)

	logrus.WithFields(f).Info(msg)
}

func Warn(msg string, params ...Parameter) {
	f := transform(params)

	logrus.WithFields(f).Warning(msg)
}

func Error(msg string, params ...Parameter) {
	f := transform(params)

	logrus.WithFields(f).Error(msg)
}

func Fatal(msg string, params ...Parameter) {
	f := transform(params)

	logrus.WithFields(f).Fatal(msg)
}

func Panic(msg string, params ...Parameter) {
	f := transform(params)

	logrus.WithFields(f).Panic(msg)
}
