package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAllSet(t *testing.T) {
	actual, _ := validateAllSet(TestStruct{String: "a", Number: 1, Array: []string{}})
	assert.True(t, actual)
}

func TestValidateNoStringSet(t *testing.T) {
	valid, missing := validateAllSet(TestStruct{Number: 1, Array: []string{}})
	assert.False(t, valid)
	assert.Equal(t, &[]string{"String"}, missing)
}

func TestValidateNoIntSet(t *testing.T) {
	valid, missing := validateAllSet(TestStruct{String: "a", Array: []string{}})
	assert.False(t, valid)
	assert.Equal(t, &[]string{"Number"}, missing)
}

func TestValidateNoArraySet(t *testing.T) {
	valid, missing := validateAllSet(TestStruct{String: "a", Number: 1})
	assert.False(t, valid)
	assert.Equal(t, &[]string{"Array"}, missing)
}

func TestValidateNothingSet(t *testing.T) {
	valid, missing := validateAllSet(TestStruct{})
	assert.False(t, valid)
	assert.Equal(t, &[]string{"String", "Number", "Array"}, missing)
}

func TestValidateEmail(t *testing.T) {
	valid := validateEmail("valid123.!#$%&'*+-/=?^_`{|}~@d.omain.tld")
	assert.True(t, valid)
}

func TestValidateEmailOnInvalidChar(t *testing.T) {
	valid := validateEmail("invalid\"@domain.tld")
	assert.False(t, valid)
}

func TestValidateEmailOnMissingRecipient(t *testing.T) {
	valid := validateEmail("@domain.tld")
	assert.False(t, valid)
}

func TestValidateEmailOnMissingDomain(t *testing.T) {
	valid := validateEmail("invalid")
	assert.False(t, valid)
}

func TestValidateEmailOnMissingTld(t *testing.T) {
	valid := validateEmail("invalid@domain")
	assert.False(t, valid)
}

func TestValidatePermissions(t *testing.T) {
	valid := validatePermissions([]string{"ADMIN", "CRT_LIST", "MOD_LIST", "CRT_USER", "MOD_USER"})
	assert.True(t, valid)
}

func TestValidatePermissionsIncompleteList(t *testing.T) {
	valid := validatePermissions([]string{"ADMIN"})
	assert.True(t, valid)
}

func TestValidatePermissionsDuplicate(t *testing.T) {
	valid := validatePermissions([]string{"ADMIN", "ADMIN"})
	assert.False(t, valid)
}

func TestValidatePermissionsInvalid(t *testing.T) {
	valid := validatePermissions([]string{"INVALID"})
	assert.False(t, valid)
}

type TestStruct struct {
	String string
	Number int
	Array  []string
}
