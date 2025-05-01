package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

var creds = &model.LoginRequest{Email: "example@domain.tld", Password: "pass"}
var hashBytes, _ = bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
var hash = string(hashBytes)

func TestLogin(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPassword", creds.Email).Return(1, hash)
	man := AuthManager{db: dbMock}

	actualToken, actualId, err := man.Login(creds)

	dbMock.AssertExpectations(t)
	require.Nil(t, err)

	token, e := jwt.Parse(actualToken, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and not something else
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return man.jwtSecret, nil
	})
	require.NoError(t, e)
	assert.True(t, token.Valid)
	assert.Equal(t, 1, actualId)
	assert.Equal(t, 1, int(token.Claims.(jwt.MapClaims)["sub"].(float64)))
}

func TestLoginReturnsNoUserError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.ErrNoRows)
	dbMock.On("GetUserPassword", defaultUserRequest.Email).Return(1, hash)
	man := AuthManager{db: dbMock}
	actualToken, actualId, err := man.Login(creds)

	assert.Equal(t, util.NewNoUserEmailError(defaultUserRequest.Email, nil), err)
	assert.Empty(t, actualToken)
	assert.Empty(t, actualId)
}

func TestLoginReturnsGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("GetUserPassword", defaultUserRequest.Email).Return(1, hash)
	man := AuthManager{db: dbMock}
	actualToken, actualId, err := man.Login(creds)

	assert.Equal(t, util.NewGenericError(nil), err)
	assert.Empty(t, actualToken)
	assert.Empty(t, actualId)
}

func TestLoginReturnsErrorFromIncorrectPassword(t *testing.T) {
	creds := &model.LoginRequest{Email: "example@domain.tld", Password: "passs"}
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPassword", creds.Email).Return(1, hash)
	man := AuthManager{db: dbMock}
	_, _, err := man.Login(creds)

	assert.Equal(t, util.ErrIncorrectPassword, err.Code)
	assert.Equal(t, "Incorrect password for user 'example@domain.tld'", err.Message)
}

func TestCheckToken(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UserExists", 1).Return(true)
	man := AuthManager{db: dbMock}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 1,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	signedToken, err := token.SignedString(man.jwtSecret)
	require.NoError(t, err)

	id, e := man.CheckTokenValidity(signedToken)

	dbMock.AssertExpectations(t)
	assert.Nil(t, e)
	assert.Equal(t, 1, id)
}

func TestCheckTokenReturnsErrorOnInvalidToken(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UserExists", mock.Anything, mock.Anything).Panic("UserExists should not have been called")
	man := AuthManager{db: dbMock}
	_, err := man.CheckTokenValidity("invalid")

	assert.Equal(t, util.ErrInvalidToken, err.Code)
	assert.Equal(t, "Invalid token provided", err.Message)
}

func TestCheckTokenReturnsErrorOnInvalidSub(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("UserExists", 1).Return(false)
	man := AuthManager{db: dbMock}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "invalid",
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	signedToken, err := token.SignedString(man.jwtSecret)
	require.NoError(t, err)

	_, e := man.CheckTokenValidity(signedToken)

	assert.Equal(t, util.NewInvalidTokenError(nil), e)
}

func TestCheckTokenReturnsFalseOnNoUserError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.ErrNoRows)
	dbMock.On("UserExists", 1).Return(false)
	man := AuthManager{db: dbMock}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 1,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	signedToken, err := token.SignedString(man.jwtSecret)
	require.NoError(t, err)

	_, e := man.CheckTokenValidity(signedToken)
	assert.Equal(t, util.NewInvalidTokenError(nil), e)
}

func TestCheckTokenReturnsGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("UserExists", 1).Return(false)
	man := AuthManager{db: dbMock}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 1,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	signedToken, err := token.SignedString(man.jwtSecret)
	require.NoError(t, err)

	_, e := man.CheckTokenValidity(signedToken)

	assert.Equal(t, util.NewInvalidTokenError(nil), e)
}

func TestCheckPerms(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"CRT_USER", "CRT_LIST"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckPerms(1, "CRT_LIST")

	require.Nil(t, err)
	assert.True(t, valid)
}

func TestCheckPermsTrueOnAdmin(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"ADMIN"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckPerms(1, "CRT_LIST")

	require.Nil(t, err)
	assert.True(t, valid)
}

func TestCheckPermsWhenFalse(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"CRT_USER"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckPerms(1, "CRT_LIST")

	assert.Equal(t, util.NewNoPermissionError("CRT_LIST", nil), err)
	assert.False(t, valid)
}

func TestCheckPermsReturnsGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("GetUserPerms", 1).Return([]string{})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckPerms(1, "CRT_LIST")

	assert.Equal(t, util.NewGenericError(nil), err)
	assert.False(t, valid)
}

func TestCheckUserPerms(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"CRT_USER", "CRT_LIST"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckUserPerms(1, "CRT_USER", []string{"CRT_USER", "CRT_LIST"})

	require.Nil(t, err)
	assert.True(t, valid)
}

func TestCheckUserPermsTrueOnAdmin(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"ADMIN"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckUserPerms(1, "CRT_USER", []string{"CRT_USER", "CRT_LIST"})

	require.Nil(t, err)
	assert.True(t, valid)
}

func TestCheckUserPermsTrueOnLowercase(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"CRT_USER", "CRT_LIST"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckUserPerms(1, "CRT_USER", []string{"crt_user", "crt_list"})

	require.Nil(t, err)
	assert.True(t, valid)
}

func TestCheckUserPermsWhenNoActionPerm(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"CRT_LIST"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckUserPerms(1, "CRT_USER", []string{"CRT_LIST"})

	assert.Equal(t, util.NewNoPermissionError("CRT_USER", nil), err)
	assert.False(t, valid)
}

func TestCheckUserPermsWhenNoRequiredPerm(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"CRT_USER"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckUserPerms(1, "CRT_USER", []string{"CRT_LIST"})

	assert.Equal(t, util.NewNoPermissionError("CRT_LIST", nil), err)
	assert.False(t, valid)
}

func TestCheckUserPermsWhenInvalidPerm(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPerms", 1).Return([]string{"CRT_USER"})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckUserPerms(1, "CRT_USER", []string{"INVALID_PERM"})

	msg := "Missing or duplicate permission. Valid permissions- ADMIN,CRT_LIST,MOD_LIST,CRT_USER,MOD_USER"
	assert.Equal(t, util.NewInvalidObjectError(msg, nil), err)
	assert.False(t, valid)
}

func TestCheckUserPermsReturnsGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("GetUserPerms", 1).Return([]string{})
	man := AuthManager{db: dbMock}

	valid, err := man.CheckUserPerms(1, "CRT_USER", []string{"CRT_LIST"})

	assert.Equal(t, util.NewGenericError(nil), err)
	assert.False(t, valid)
}

func TestCheckListMods(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPermsAndModStatus", 1, 1).Return([]string{"CRT_USER", "MOD_LIST"}, false)
	man := AuthManager{db: dbMock}

	valid, err := man.CheckListMods(1, 1)

	require.Nil(t, err)
	assert.True(t, valid)
}

func TestCheckListModsOnAdmin(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPermsAndModStatus", 1, 1).Return([]string{"ADMIN"}, false)
	man := AuthManager{db: dbMock}

	valid, err := man.CheckListMods(1, 1)

	require.Nil(t, err)
	assert.True(t, valid)
}

func TestCheckListModsOnMod(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPermsAndModStatus", 1, 1).Return([]string{}, true)
	man := AuthManager{db: dbMock}

	valid, err := man.CheckListMods(1, 1)

	require.Nil(t, err)
	assert.True(t, valid)
}

func TestCheckListModsWhenFalse(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetUserPermsAndModStatus", 1, 1).Return([]string{"CRT_USER"}, false)
	man := AuthManager{db: dbMock}

	valid, err := man.CheckListMods(1, 1)

	assert.Equal(t, util.NewNoPermissionError("MOD_LIST", nil), err)
	assert.False(t, valid)
}

func TestCheckListModsReturnsGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("GetUserPermsAndModStatus", 1, 1).Return([]string{}, false)
	man := AuthManager{db: dbMock}

	valid, err := man.CheckListMods(1, 1)

	assert.Equal(t, util.NewGenericError(nil), err)
	assert.False(t, valid)
}
