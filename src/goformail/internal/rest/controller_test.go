package rest

import (
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"net/http"
	"testing"
)

func TestNewController(t *testing.T) {
	actual := NewController(&service.ListManager{}, &service.UserManager{}, &service.AuthManager{}, &db.Db{})
	expected := &Controller{
		&service.ListManager{},
		&service.UserManager{},
		&service.AuthManager{},
		&db.Db{},
		http.DefaultServeMux,
	}
	assert.Equal(t, expected, actual)
}
