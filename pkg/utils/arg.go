package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrFlagTagNotFound = errors.New("arg: given field doesn't have a `flag` tag")
)

// Uses reflection to retrieve the `flag` tag of a field.
// The value of the `flag` field with the value of the field is
// used to construct a POSIX long flag argument string.
func GetLongFlagFormOfField(fieldValue reflect.Value, fieldType reflect.StructField) (string, error) {
	flagTag := fieldType.Tag.Get("flag")
	if flagTag == "" {
		return "", ErrFlagTagNotFound
	}

	switch fieldValue.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("--%v", flagTag), nil
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		args := make([]string, 0)
		for i := 0; i < fieldValue.Len(); i++ {
			args = append(args, fieldValue.Index(i).String())
		}

		return fmt.Sprintf("--%v=%v", flagTag, strings.Join(args, ",")), nil
	default:
		return fmt.Sprintf("--%v=%v", flagTag, fieldValue.String()), nil
	}
}

// Uses reflection to transform a struct containing fields with `flag` tags
// to a string slice of POSIX compliant long form arguments.
func GetArgumentFormOfStruct(strt interface{}) (flags []string) {
	numberOfFields := reflect.ValueOf(strt).NumField()
	for i := 0; i < numberOfFields; i++ {
		fieldValue := reflect.ValueOf(strt).Field(i)
		fieldType := reflect.TypeOf(strt).Field(i)

		if fieldValue.String() == "" {
			continue
		}

		flagFormOfField, err := GetLongFlagFormOfField(fieldValue, fieldType)
		if err != nil {
			continue
		}

		flags = append(flags, flagFormOfField)
	}

	return
}
