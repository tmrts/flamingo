// Package logfield contains log field implementations
package logfield

import (
	"fmt"

	"github.com/tmrts/flamingo/pkg/flog"
)

func Error(err interface{}) flog.Field {
	return func() (string, string) {
		return "error", err.(error).Error()
	}
}

func Event(event interface{}, data ...interface{}) flog.Field {
	return func() (string, string) {
		return "event", fmt.Sprintf(event.(string), data...)
	}
}
