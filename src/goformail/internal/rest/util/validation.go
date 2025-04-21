package util

import (
	"reflect"
)

func ValidateAllSet(object interface{}) bool {
	if reflect.TypeOf(object).Kind() != reflect.Struct {
		return false
	}

	allSet := true
	value := reflect.ValueOf(object)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		allSet = field.IsValid() && !field.IsZero()
		if !allSet {
			break
		}
	}
	return allSet
}
