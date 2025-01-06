package userImpl

import (
	"errors"
	"fmt"
	"go-sns/database/dataAccess/interfaces"
	"go-sns/database/mocks"
	"go-sns/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_Success(t *testing.T) {
	mockDB := new(mocks.Database)
	mockDB.On(
		"PrepareAndExecute",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	mockDB.On("GetLastInsertedId").Return(1, nil)

	now := time.Now().Truncate(time.Second)
	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)

	timestamp := models.NewDateTimeStamp(now, now)
	user := models.NewUser(-1, "Test User", "Test Email", *timestamp)
	password := "Abcd1234!"
	err := dao.Create(user, password)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestCreate_InvalidID(t *testing.T) {
	mockDB := new(mocks.Database)
	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)

	invalidID := 1
	user := models.NewUser(invalidID, "Test User", "Test Email", models.DateTimeStamp{})

	password := "Abcd1234!"
	err := dao.Create(user, password)

	expectedError := "action=UserDAOImpl.Create msg=Cannot create a user data with an existing ID. id: " + string(rune(invalidID))
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	mockDB.AssertNotCalled(t, "PrepareAndExecute")
}

func TestCreate_QueryExecutionFailure(t *testing.T) {
	mockDB := new(mocks.Database)

	queryError := errors.New("database error")
	mockDB.On(
		"PrepareAndExecute",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(queryError)

	now := time.Now().Truncate(time.Second)
	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)

	timestamp := models.NewDateTimeStamp(now, now)
	user := models.NewUser(-1, "Test User", "Test Email", *timestamp)
	password := "Abcd1234!"
	err := dao.Create(user, password)

	expectedError := "action=UserDAOImpl.Create msg=Error executing query: database error"
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	mockDB.AssertNotCalled(t, "GetLastInsertedId")
	mockDB.AssertExpectations(t)
}

func TestCreate_FetchingIdQueryFailure(t *testing.T) {
	mockDB := new(mocks.Database)

	mockDB.On(
		"PrepareAndExecute",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	fetchingError := errors.New("database error")

	mockDB.On("GetLastInsertedId").Return(-1, fetchingError)

	now := time.Now().Truncate(time.Second)
	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)

	timestamp := models.NewDateTimeStamp(now, now)
	user := models.NewUser(-1, "Test User", "Test Email", *timestamp)
	password := "Abcd1234!"
	err := dao.Create(user, password)

	expectedError := "action=UserDAOImpl.Create msg=Error fetching last insert ID: database error"
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	mockDB.AssertExpectations(t)
}



func generateMockUsers(count int, now time.Time) []map[string]interface{} {
	mockData := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		mockData[i] = map[string]interface{}{
			"id":          i + 1,
			"name":       fmt.Sprintf("Test name %d", i), 
			"email": fmt.Sprintf("Test email %d", i),
			"password": "#####",
			"created_at":  now,
			"updated_at":  now,
		}
	}
	return mockData
}

func TestGetAll_Success(t *testing.T){
	var recordSize int = 100
    now := time.Now().Truncate(time.Second)

	mockDB := new(mocks.Database)
	mockDB.On(
		"GetTableLength",
		"users",).Return(recordSize, nil)

	mockDB.On(
		"PrepareAndFetchAll",
		mock.Anything,
		recordSize).Return(generateMockUsers(recordSize, now), nil)

	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)


	users, err := dao.GetAll()

	assert.NoError(t, err)
	assert.Len(t, users, recordSize)
	nameId := 10
	emailId := 50
	expectedName := fmt.Sprintf("Test name %v", nameId)
	expecteEmail := fmt.Sprintf("Test email %v",emailId)

	assert.Equal(t, expectedName, users[nameId].GetName())
	assert.Equal(t, expecteEmail,users[emailId].GetEmail())

	mockDB.AssertExpectations(t)
}

func TestGetAllWithLimit_Success(t *testing.T){
	var limit int = 10
	var recordSize int = 100
	
    now := time.Now().Truncate(time.Second)

	mockDB := new(mocks.Database)

	mockDB.On(
		"GetTableLength",
		"users",).Return(recordSize, nil)

	mockDB.On(
		"PrepareAndFetchAll",
		mock.Anything,
		limit).Return(generateMockUsers(limit,now), nil)

	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)

	users, err := dao.GetAll(limit)

	assert.NoError(t, err)
	assert.Len(t, users, limit)

	nameId := 5
	emailId := 3
	expectedName := fmt.Sprintf("Test name %v", nameId)
	expecteEmail := fmt.Sprintf("Test email %v",emailId)

	assert.Equal(t, expectedName, users[nameId].GetName())
	assert.Equal(t, expecteEmail,users[emailId].GetEmail())

	mockDB.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "GetTableLength")
}


func TestGetAllWithLimit_OverSize(t *testing.T){
	var recordSize int = 100
	var limit int = 1e9

    now := time.Now().Truncate(time.Second)

	mockDB := new(mocks.Database)

	mockDB.On(
		"GetTableLength",
		"users",).Return(recordSize, nil)

	mockDB.On(
		"PrepareAndFetchAll",
		mock.Anything,
		recordSize).Return(generateMockUsers(recordSize,now), nil)

	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)

	users, err := dao.GetAll(limit)

	assert.NoError(t, err)
	assert.Len(t, users, recordSize)
	nameId := 10
	emailId := 50
	expectedName := fmt.Sprintf("Test name %v", nameId)
	expecteEmail := fmt.Sprintf("Test email %v",emailId)

	assert.Equal(t, expectedName, users[nameId].GetName())
	assert.Equal(t, expecteEmail,users[emailId].GetEmail())

	mockDB.AssertExpectations(t)
}


func TestGetById_Success(t *testing.T) {
	var id int = 10
	expectedName := fmt.Sprintf("Test name %v", id)
	now := time.Now().Truncate(time.Second)
	mockDB := new(mocks.Database)
	mockDB.On(
		"PrepareAndFetchAll",
		mock.Anything,
		id).Return([]map[string]interface{}{
		{
			"id":         id,
			"name":       expectedName,
			"email":      "",
			"password":   "",
			"created_at": now,
			"updated_at": now,
		},
	}, nil)

	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)

	user, err := dao.GetById(id)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, id, user.GetId())
	assert.Equal(t, expectedName, user.GetName())

	mockDB.AssertExpectations(t)
}

func TestGetById_NotFound(t *testing.T) {
	var invalidId int = -1

	mockDB := new(mocks.Database)
	mockDB.On("PrepareAndFetchAll", mock.Anything, invalidId).
		Return([]map[string]interface{}{}, nil)

	var dao interfaces.UserDAO = NewUserDAOImpl(mockDB)
	user, err := dao.GetById(invalidId)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, fmt.Sprintf("no user found with id %d", invalidId))
	mockDB.AssertExpectations(t)
}
