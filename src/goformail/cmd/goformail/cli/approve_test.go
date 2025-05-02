package cli

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"testing"
)

func TestApproveEmail(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("SetEmailAsApproved", 1).Return(nil)
	approveEmail([]string{"1"}, dbMock)
}
