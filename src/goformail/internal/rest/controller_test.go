package rest

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/test/mock"
	"net/http"
	"reflect"
	"testing"
)

func TestNewController(t *testing.T) {
	actual := NewController(mock.Configs, nil)
	expected := &Controller{mock.Configs, nil, http.DefaultServeMux}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("NewController created invalid controller. Expected: '%v', got '%v'", expected, actual)
	}
}
