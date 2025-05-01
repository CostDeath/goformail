package service

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"log"
	"reflect"
	"regexp"
	"strings"
)

func validateAllSet(object interface{}) (bool, *[]string) {
	objType := reflect.TypeOf(object)
	if objType.Kind() != reflect.Struct {
		return false, nil
	}

	allSet := true
	var missing []string
	objValue := reflect.ValueOf(object)
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		set := field.IsValid() && !field.IsZero()
		if !set && field.Kind() != reflect.Bool {
			allSet = false
			missing = append(missing, objType.Field(i).Name)
		}
	}
	return allSet, &missing
}

func validatePermissionsSet(user model.UserRequest) bool {
	field := reflect.ValueOf(user).FieldByName("Permissions")
	return field.IsValid() && !field.IsZero()
}

func validateListPropsSet(list model.ListRequest) *model.ListOverrides {
	rcpt := reflect.ValueOf(list).FieldByName("Recipients")
	mods := reflect.ValueOf(list).FieldByName("Mods")
	senders := reflect.ValueOf(list).FieldByName("ApprovedSenders")

	return &model.ListOverrides{
		Recipients:      rcpt.IsValid() && !rcpt.IsZero(),
		Mods:            mods.IsValid() && !mods.IsZero(),
		ApprovedSenders: senders.IsValid() && !senders.IsZero(),
	}
}

func validateEmail(txt string) bool {
	matches, err := regexp.Match(
		`^([A-z0-9+.!#$%&'*+-/=?^_{|}~][-A-z0-9+.!#$%&'*+-/=?^_{|}~]*)@(([a-z0-9][-a-z0-9]*\.)([-a-z0-9]+\.)*[a-z]{2,})$`,
		[]byte(txt),
	)
	if err != nil {
		log.Print(err)
	}
	return matches
}

func validatePermissions(perms []string) bool {
	seen := make(map[string]bool)
	allowed := make(map[string]bool)
	for _, perm := range model.Permissions {
		allowed[perm] = true
	}

	for _, perm := range perms {
		perm := strings.ToUpper(perm)
		if !allowed[perm] || seen[perm] {
			return false
		}
		seen[perm] = true
	}
	return true
}
