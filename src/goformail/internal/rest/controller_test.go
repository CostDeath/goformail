package rest

import (
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"net/http"
	"testing"
)

func TestNewController(t *testing.T) {
	actual := NewController(util.MockConfigs, &db.Db{}, &service.UserManager{}, &service.AuthManager{})
	expected := &Controller{
		util.MockConfigs,
		&db.Db{},
		&service.UserManager{},
		&service.AuthManager{},
		http.DefaultServeMux,
	}
	assert.Equal(t, expected, actual)
}
