package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type IAuthManager interface {
	Login(creds *model.LoginRequest) (string, *util.Error)
	CheckTokenValidity(tokenStr string) (int, *util.Error)
	CheckPerms(id int, required string) (bool, *util.Error)
	CheckUserPerms(id int, action string, required []string) (bool, *util.Error)
}

type AuthManager struct {
	IAuthManager
	db        db.IDb
	jwtSecret []byte
}

func NewAuthManager(db db.IDb, jwtSecret *[]byte) *AuthManager {
	return &AuthManager{db: db, jwtSecret: *jwtSecret}
}

func (man *AuthManager) Login(creds *model.LoginRequest) (string, *util.Error) {
	// Get hash from db
	id, hash, dbErr := man.db.GetUserPassword(creds.Email)
	if dbErr != nil {
		switch dbErr.Code {
		case db.ErrNoRows:
			return "", util.NewNoUserEmailError(creds.Email, dbErr.Err)
		default:
			return "", util.NewGenericError(dbErr.Err)
		}
	}

	// Check password matches hash in db
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", util.NewIncorrectPasswordError(creds.Email, err)
		}
		return "", util.NewEncryptionError(err)
	}

	// Generate signed JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	signedToken, err := token.SignedString(man.jwtSecret)
	if err != nil {
		return "", util.NewEncryptionError(err)
	}

	return signedToken, nil
}

func (man *AuthManager) CheckTokenValidity(tokenStr string) (int, *util.Error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and not something else
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return man.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, util.NewInvalidTokenError(err)
	}

	// Token is valid — get user's id
	claims := token.Claims.(jwt.MapClaims)
	floatId, ok := claims["sub"].(float64)
	if !ok {
		return 0, util.NewInvalidTokenError(nil)
	}
	id := int(floatId)

	if exists, err := man.db.UserExists(id); err != nil || !exists {
		return 0, util.NewInvalidTokenError(nil)
	}
	return id, nil
}

func (man *AuthManager) CheckPerms(id int, required string) (bool, *util.Error) {
	perms, dbErr := man.db.GetUserPerms(id)
	if dbErr != nil {
		return false, util.NewGenericError(dbErr.Err)
	}

	for _, perm := range perms {
		if perm == "ADMIN" || perm == required {
			return true, nil
		}
	}

	return false, util.NewNoPermissionError(required, nil)
}

func (man *AuthManager) CheckUserPerms(id int, action string, required []string) (bool, *util.Error) {
	if !validatePermissions(required) {
		msg := "Missing or duplicate permission. Valid permissions- " + strings.Join(model.Permissions, ",")
		return false, util.NewInvalidObjectError(msg, nil)
	}

	perms, dbErr := man.db.GetUserPerms(id)
	if dbErr != nil {
		return false, util.NewGenericError(dbErr.Err)
	}

	seen := make(map[string]bool)
	for _, perm := range perms {
		if perm == "ADMIN" {
			return true, nil
		}
		seen[perm] = true
	}

	if !seen[action] {
		return false, util.NewNoPermissionError(action, nil)
	}

	for _, perm := range required {
		perm := strings.ToUpper(perm)
		if !seen[perm] {
			return false, util.NewNoPermissionError(perm, nil)
		}
	}

	return true, nil
}
