package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"testing"
)

var defaultListRequest = &model.ListRequest{
	Name:            "test",
	Recipients:      []string{"rcpt1@domain.tld", "rcpt2@domain.tld"},
	Mods:            []int64{1, 2},
	ApprovedSenders: []string{"sdr1@domain.tld", "sdr2@domain.tld"},
}
var defaultListRequestToFix = &model.ListRequest{
	Name:            defaultListRequest.Name,
	Recipients:      []string{"RCPT1@domain.tld", "RCPT1@domain.tld", "RCPT2@domain.tld"},
	Mods:            []int64{1, 2, 3},
	ApprovedSenders: []string{"SDR1@domain.tld", "SDR2@domain.tld", "SDR2@domain.tld"},
}
var defaultListResponse = &model.ListResponse{
	Id:              1,
	Name:            defaultListRequest.Name,
	Recipients:      defaultListRequest.Recipients,
	Mods:            defaultListRequest.Mods,
	ApprovedSenders: defaultListRequest.ApprovedSenders,
}

func TestGetList(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("GetList", 1).Return(defaultListResponse)
	man := ListManager{db: mockObj}
	actual, err := man.GetList(1)

	mockObj.AssertExpectations(t)
	require.Nil(t, err)
	assert.Equal(t, defaultListResponse, actual)
}

func TestGetListReturnsNoListError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("GetList", 1).Return(defaultListResponse)
	man := ListManager{db: mockObj}
	actual, err := man.GetList(1)

	mockObj.AssertExpectations(t)
	assert.Equal(t, util.NewNoListError(1, nil), err)
	assert.Nil(t, actual)
}

func TestGetListReturnsGenericError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("GetList", 1).Return(defaultListResponse)
	man := ListManager{db: mockObj}
	actual, err := man.GetList(1)

	mockObj.AssertExpectations(t)
	assert.Equal(t, util.NewGenericError(nil), err)
	assert.Nil(t, actual)
}

func TestCreateList(t *testing.T) {
	dbMock := new(db.IDbMock)
	valid := defaultListRequest.Mods
	dbMock.On("UsersExist", defaultListRequestToFix.Mods).Return(valid)
	dbMock.On("CreateList", defaultListRequest).Return(1)
	man := ListManager{db: dbMock}

	id, err := man.CreateList(defaultListRequestToFix)

	dbMock.AssertExpectations(t)
	require.Nil(t, err)
	assert.Equal(t, 1, id)
}

func TestCreateListReturnsInvalidObjectErrorOnMissingFields(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UsersExist", mock.Anything).Panic("UsersExist should not have been called")
	dbMock.On("CreateList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	id, err := man.CreateList(&model.ListRequest{})

	expected := util.NewInvalidObjectError("Missing field(s) in list: Name, Recipients, Mods, ApprovedSenders", nil)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, id)
}

func TestCreateListReturnsInvalidObjectErrorOnInvalidName(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UsersExist", mock.Anything).Panic("UsersExist should not have been called")
	dbMock.On("CreateList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	id, err := man.CreateList(&model.ListRequest{
		Name:            "invalid@",
		Recipients:      defaultListRequest.Recipients,
		Mods:            defaultListRequest.Mods,
		ApprovedSenders: defaultListRequest.ApprovedSenders})

	expected := util.NewInvalidObjectError("Invalid list name 'invalid@' (must not include domain)", nil)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, id)
}

func TestCreateListReturnsInvalidObjectErrorOnInvalidRecipient(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UsersExist", mock.Anything).Panic("UsersExist should not have been called")
	dbMock.On("CreateList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	id, err := man.CreateList(&model.ListRequest{
		Name:            defaultListRequest.Name,
		Recipients:      []string{"invalid"},
		Mods:            defaultListRequest.Mods,
		ApprovedSenders: defaultListRequest.ApprovedSenders})

	expected := util.NewInvalidObjectError("Invalid recipient email address 'invalid'", nil)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, id)
}

func TestCreateListReturnsInvalidObjectErrorOnInvalidSender(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UsersExist", mock.Anything).Panic("UsersExist should not have been called")
	dbMock.On("CreateList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	id, err := man.CreateList(&model.ListRequest{
		Name:            defaultListRequest.Name,
		Recipients:      defaultListRequest.Recipients,
		Mods:            defaultListRequest.Mods,
		ApprovedSenders: []string{"invalid"}})

	expected := util.NewInvalidObjectError("Invalid sender email address 'invalid'", nil)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, id)
}

func TestCreateListReturnsGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("UsersExist", defaultListRequest.Mods).Return(defaultListResponse.Mods)
	dbMock.On("CreateList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	id, err := man.CreateList(defaultListRequest)

	assert.Equal(t, util.NewGenericError(nil), err)
	assert.Equal(t, 0, id)
}

func TestUpdateList(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UsersExist", defaultListRequestToFix.Mods).Return(defaultListRequest.Mods)
	dbMock.On("PatchList", 1, defaultListRequest, &model.ListOverrides{
		Recipients: true, Mods: true, ApprovedSenders: true,
	}).Return()
	man := ListManager{db: dbMock}

	err := man.UpdateList(1, defaultListRequestToFix)

	dbMock.AssertExpectations(t)
	require.Nil(t, err)
}

func TestUpdateListDoesNotOverrideWhenNullProps(t *testing.T) {
	dbMock := new(db.IDbMock)
	var missing []int64
	list := &model.ListRequest{Name: "name"}
	dbMock.On("UsersExist", missing).Return(missing)
	dbMock.On("PatchList", 1, list, &model.ListOverrides{}).Return()
	man := ListManager{db: dbMock}

	err := man.UpdateList(1, list)

	dbMock.AssertExpectations(t)
	require.Nil(t, err)
}

func TestUpdateListReturnsInvalidObjectErrorOnInvalidName(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UsersExist", mock.Anything).Panic("UsersExist should not have been called")
	dbMock.On("PatchList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	err := man.UpdateList(1, &model.ListRequest{
		Name:            "invalid@",
		Recipients:      defaultListRequest.Recipients,
		Mods:            defaultListRequest.Mods,
		ApprovedSenders: defaultListRequest.ApprovedSenders})

	expected := util.NewInvalidObjectError("Invalid list name 'invalid@' (must not include domain)", nil)
	assert.Equal(t, expected, err)
}

func TestUpdateListReturnsInvalidObjectErrorOnInvalidRecipient(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UsersExist", mock.Anything).Panic("UsersExist should not have been called")
	dbMock.On("PatchList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	err := man.UpdateList(1, &model.ListRequest{
		Name:            defaultListRequest.Name,
		Recipients:      []string{"invalid"},
		Mods:            defaultListRequest.Mods,
		ApprovedSenders: defaultListRequest.ApprovedSenders})

	expected := util.NewInvalidObjectError("Invalid recipient email address 'invalid'", nil)
	assert.Equal(t, expected, err)
}

func TestUpdateListReturnsInvalidObjectErrorOnInvalidSender(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("UsersExist", mock.Anything).Panic("UsersExist should not have been called")
	dbMock.On("PatchList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	err := man.UpdateList(1, &model.ListRequest{
		Name:            defaultListRequest.Name,
		Recipients:      defaultListRequest.Recipients,
		Mods:            defaultListRequest.Mods,
		ApprovedSenders: []string{"invalid"}})

	expected := util.NewInvalidObjectError("Invalid sender email address 'invalid'", nil)
	assert.Equal(t, expected, err)
}

func TestUpdateListReturnsGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("UsersExist", defaultListRequest.Mods).Return(defaultListResponse.Mods)
	dbMock.On("PatchList", mock.Anything).Panic("CreateList should not have been called")
	man := ListManager{db: dbMock}

	err := man.UpdateList(1, defaultListRequest)

	assert.Equal(t, util.NewGenericError(nil), err)
}

func TestDeleteList(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("DeleteList", 1).Return()
	man := ListManager{db: mockObj}
	err := man.DeleteList(1)

	mockObj.AssertExpectations(t)
	require.Nil(t, err)
}

func TestDeleteListReturnsNoListError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("DeleteList", 1).Return()
	man := ListManager{db: mockObj}
	err := man.DeleteList(1)

	mockObj.AssertExpectations(t)
	assert.Equal(t, util.NewNoListError(1, nil), err)
}

func TestDeleteListReturnsGenericError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("DeleteList", 1).Return()
	man := ListManager{db: mockObj}
	err := man.DeleteList(1)

	assert.Equal(t, util.NewGenericError(nil), err)
}

func TestGetAllLists(t *testing.T) {
	mockObj := new(db.IDbMock)
	expected := &[]*model.ListResponse{defaultListResponse}
	mockObj.On("GetAllLists").Return(expected)
	man := ListManager{db: mockObj}
	actual, err := man.GetAllLists()

	mockObj.AssertExpectations(t)
	require.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetListsReturnsGenericError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("GetAllLists").Return(&[]*model.ListResponse{defaultListResponse})
	man := ListManager{db: mockObj}
	actual, err := man.GetAllLists()

	assert.Equal(t, util.NewGenericError(nil), err)
	assert.Nil(t, actual)
}
