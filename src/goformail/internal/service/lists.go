package service

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"strings"
)

type IListManager interface {
	GetList(id int) (*model.ListResponse, *util.Error)
	CreateList(list *model.ListRequest) (int, *util.Error)
	UpdateList(id int, list *model.ListRequest, hasLocked bool) *util.Error
	DeleteList(id int) *util.Error
	GetAllLists() (*[]*model.ListResponse, *util.Error)
}

type ListManager struct {
	IListManager
	db db.IDb
}

func NewListManager(db *db.Db) *ListManager {
	return &ListManager{db: db}
}

func (man *ListManager) GetList(id int) (*model.ListResponse, *util.Error) {
	list, err := man.db.GetList(id)
	if err != nil {
		switch err.Code {
		case db.ErrNoRows:
			return nil, util.NewNoListError(id, err.Err)
		default:
			return nil, util.NewGenericError(err.Err)
		}
	}

	return list, nil
}

func (man *ListManager) CreateList(list *model.ListRequest) (int, *util.Error) {
	// Ensure payload is valid
	if valid, missing := validateAllSet(*list); !valid {
		return 0, util.NewInvalidObjectError("Missing field(s) in list: "+strings.Join(*missing, ", "), nil)
	}

	list.Name = strings.ToLower(list.Name) // want to store lowercase, to prevent duplicates
	if !validateEmail(list.Name + "@domain.tld") {
		return 0, util.NewInvalidObjectError("Invalid list name '"+list.Name+"' (must not include domain)", nil)
	}

	seen := make(map[string]bool)
	var newRcptList []string
	for _, recipient := range list.Recipients {
		recipient = strings.ToLower(recipient) // want to store lowercase, to prevent duplicates
		if !validateEmail(recipient) {
			return 0, util.NewInvalidObjectError("Invalid recipient email address '"+recipient+"'", nil)
		}
		if !seen[recipient] {
			newRcptList = append(newRcptList, recipient)
		}
		seen[recipient] = true
	}
	list.Recipients = newRcptList

	seen = make(map[string]bool)
	var newSdrList []string
	for _, sender := range list.ApprovedSenders {
		sender = strings.ToLower(sender) // want to store lowercase, to prevent duplicates
		if !validateEmail(sender) {
			return 0, util.NewInvalidObjectError("Invalid sender email address '"+sender+"'", nil)
		}
		if !seen[sender] {
			newSdrList = append(newSdrList, sender)
		}
		seen[sender] = true
	}
	list.ApprovedSenders = newSdrList

	// Verify mods exist
	if valid, err := man.db.UsersExist(list.Mods); err != nil {
		return 0, util.NewGenericError(err.Err)
	} else {
		list.Mods = valid
	}

	// Create list in db
	id, err := man.db.CreateList(list)
	if err != nil {
		switch err.Code {
		case db.ErrDuplicate:
			return 0, util.NewListAlreadyExistsError(list.Name, err.Err)
		default:
			return 0, util.NewGenericError(err.Err)
		}
	}

	return id, nil
}

func (man *ListManager) UpdateList(id int, list *model.ListRequest, hasLocked bool) *util.Error {
	list.Name = strings.ToLower(list.Name) // want to store lowercase, to prevent duplicates
	if list.Name != "" && !validateEmail(list.Name+"@domain.tld") {
		return util.NewInvalidObjectError("Invalid list name '"+list.Name+"' (must not include domain)", nil)
	}

	overrides := validateListPropsSet(*list)
	overrides.Locked = hasLocked

	seen := make(map[string]bool)
	var newRcptList []string
	for _, recipient := range list.Recipients {
		recipient = strings.ToLower(recipient) // want to store lowercase, to prevent duplicates
		if !validateEmail(recipient) {
			return util.NewInvalidObjectError("Invalid recipient email address '"+recipient+"'", nil)
		}
		if !seen[recipient] {
			newRcptList = append(newRcptList, recipient)
		}
		seen[recipient] = true
	}
	list.Recipients = newRcptList

	seen = make(map[string]bool)
	var newSdrList []string
	for _, sender := range list.ApprovedSenders {
		sender = strings.ToLower(sender) // want to store lowercase, to prevent duplicates
		if !validateEmail(sender) {
			return util.NewInvalidObjectError("Invalid sender email address '"+sender+"'", nil)
		}
		if !seen[sender] {
			newSdrList = append(newSdrList, sender)
		}
		seen[sender] = true
	}
	list.ApprovedSenders = newSdrList

	// Verify mods exist
	if valid, err := man.db.UsersExist(list.Mods); err != nil {
		return util.NewGenericError(err.Err)
	} else {
		list.Mods = valid
	}

	err := man.db.PatchList(id, list, overrides)
	if err != nil {
		switch err.Code {
		case db.ErrNoRows:
			return util.NewNoListError(id, err.Err)
		case db.ErrDuplicate:
			return util.NewListAlreadyExistsError(list.Name, err.Err)
		default:
			return util.NewGenericError(err.Err)
		}
	}

	return nil
}

func (man *ListManager) DeleteList(id int) *util.Error {
	err := man.db.DeleteList(id)
	if err != nil {
		switch err.Code {
		case db.ErrNoRows:
			return util.NewNoListError(id, err.Err)
		default:
			return util.NewGenericError(err.Err)
		}
	}

	return nil
}

func (man *ListManager) GetAllLists() (*[]*model.ListResponse, *util.Error) {
	lists, err := man.db.GetAllLists()
	if err != nil {
		return nil, util.NewGenericError(err.Err)
	}

	return lists, nil
}
