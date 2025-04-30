package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
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

func TestValidatePermissionsSet(t *testing.T) {
	var user model.UserRequest
	err := json.Unmarshal([]byte(`{"permissions": []}`), &user)
	require.NoError(t, err)

	actual := validatePermissionsSet(user)
	assert.True(t, actual)
}

func TestValidatePermissionsUnset(t *testing.T) {
	var user model.UserRequest
	err := json.Unmarshal([]byte(`{"email": "test@domain.tld"}`), &user)
	require.NoError(t, err)

	actual := validatePermissionsSet(user)
	assert.False(t, actual)
}

func TestValidateListPropsAllSet(t *testing.T) {
	var list model.ListRequest
	err := json.Unmarshal([]byte(`{"recipients": [], "mods": [], "approved_senders": []}`), &list)
	require.NoError(t, err)

	actual := validateListPropsSet(list)
	assert.Equal(t, &model.ListOverrides{Recipients: true, Mods: true, ApprovedSenders: true}, actual)
}

func TestValidateListPropsSomeSet(t *testing.T) {
	var list model.ListRequest
	err := json.Unmarshal([]byte(`{"recipients": [], "mods": []}`), &list)
	require.NoError(t, err)

	actual := validateListPropsSet(list)
	assert.Equal(t, &model.ListOverrides{Recipients: true, Mods: true}, actual)
}

func TestValidateListPropsNoneSet(t *testing.T) {
	var list model.ListRequest
	err := json.Unmarshal([]byte(`{}`), &list)
	require.NoError(t, err)

	actual := validateListPropsSet(list)
	assert.Equal(t, &model.ListOverrides{}, actual)
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
