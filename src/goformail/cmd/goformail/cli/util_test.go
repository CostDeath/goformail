package cli

import (
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"testing"
)

func TestGetListManager(t *testing.T) {
	man := getListManager(&db.Db{})
	assert.Equal(t, man, service.NewListManager(&db.Db{}))
}

func TestGetUserManager(t *testing.T) {
	man := getUserManager(&db.Db{})
	assert.Equal(t, man, service.NewUserManager(&db.Db{}))
}
