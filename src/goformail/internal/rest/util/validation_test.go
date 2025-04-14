package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAllSet(t *testing.T) {
	actual := ValidateAllSet(TestStruct{String: "a", Number: 1, Array: []string{}})
	assert.True(t, actual)
}

func TestValidateNoStringSet(t *testing.T) {
	actual := ValidateAllSet(TestStruct{Number: 1, Array: []string{}})
	assert.False(t, actual)
}

func TestValidateNoIntSet(t *testing.T) {
	actual := ValidateAllSet(TestStruct{String: "a", Array: []string{}})
	assert.False(t, actual)
}

func TestValidateNoArraySet(t *testing.T) {
	actual := ValidateAllSet(TestStruct{String: "a", Number: 1})
	assert.False(t, actual)
}

type TestStruct struct {
	String string
	Number int
	Array  []string
}
