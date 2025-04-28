package service

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"strings"
)

type IUserManager interface {
	GetUser(id int) (*model.UserResponse, *util.Error)
	CreateUser(user *model.UserRequest) (int, *util.Error)
	UpdateUser(id int, user *model.UserRequest) *util.Error
	DeleteUser(id int) *util.Error
	GetAllUsers() (*[]*model.UserResponse, *util.Error)
}

type UserManager struct {
	IUserManager
	db db.IDb
}

func NewUserManager(db *db.Db) *UserManager {
	return &UserManager{db: db}
}

func (man *UserManager) GetUser(id int) (*model.UserResponse, *util.Error) {
	user, err := man.db.GetUser(id)
	if err != nil {
		switch err.Code {
		case db.ErrNoRows:
			return nil, util.NewNoUserError(id, err.Err)
		default:
			return nil, util.NewGenericError(err.Err)
		}
	}

	return user, nil
}

func (man *UserManager) CreateUser(user *model.UserRequest) (int, *util.Error) {
	if valid, missing := validateAllSet(*user); !valid {
		return 0, util.NewInvalidObjectError("Missing field(s) in user: "+strings.Join(*missing, ", "), nil)
	}

	user.Email = strings.ToLower(user.Email) // want to store lowercase, to prevent duplicates
	if !validateEmail(user.Email) {
		return 0, util.NewInvalidObjectError("Invalid email address '"+user.Email+"'", nil)
	}

	if len(user.Permissions) != 0 && !validatePermissions(user.Permissions) {
		return 0, util.NewInvalidObjectError("Missing or duplicate permission. Valid permissions- "+
			strings.Join(model.Permissions, ", "), nil)
	}

	id, err := man.db.CreateUser(user, "hash", "salt")
	if err != nil {
		switch err.Code {
		case db.ErrDuplicate:
			return 0, util.NewUserAlreadyExistsError(user.Email, err.Err)
		default:
			return 0, util.NewGenericError(err.Err)
		}
	}

	return id, nil
}

func (man *UserManager) UpdateUser(id int, user *model.UserRequest) *util.Error {
	user.Email = strings.ToLower(user.Email) // want to store lowercase, to prevent duplicates
	if user.Email != "" && !validateEmail(user.Email) {
		return util.NewInvalidObjectError("Invalid email address '"+user.Email+"'", nil)
	}

	if len(user.Permissions) != 0 && !validatePermissions(user.Permissions) {
		return util.NewInvalidObjectError("Missing or duplicate permission. Valid permissions- "+
			strings.Join(model.Permissions, ", "), nil)
	}

	err := man.db.UpdateUser(id, user)
	if err != nil {
		switch err.Code {
		case db.ErrNoRows:
			return util.NewNoUserError(id, err.Err)
		case db.ErrDuplicate:
			return util.NewUserAlreadyExistsError(user.Email, err.Err)
		default:
			return util.NewGenericError(err.Err)
		}
	}

	return nil
}

func (man *UserManager) DeleteUser(id int) *util.Error {
	err := man.db.DeleteUser(id)
	if err != nil {
		switch err.Code {
		case db.ErrNoRows:
			return util.NewNoUserError(id, err.Err)
		default:
			return util.NewGenericError(err.Err)
		}
	}

	return nil
}

func (man *UserManager) GetAllUsers() (*[]*model.UserResponse, *util.Error) {
	users, err := man.db.GetAllUsers()
	if err != nil {
		return nil, util.NewGenericError(err.Err)
	}

	return users, nil
}
