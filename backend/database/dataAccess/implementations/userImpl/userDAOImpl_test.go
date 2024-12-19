package userImpl

import (
	"errors"
	"fmt"
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
	dao := NewUserDAOImpl(mockDB)
	timestamp := models.NewDateTimeStamp(now, now)
	user := models.NewUser(-1, "Test User", "Test Email", *timestamp)
	password := "Abcd1234!"
	err := dao.Create(user, password)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestCreate_InvalidID(t *testing.T) {
	mockDB := new(mocks.Database)
	dao := NewUserDAOImpl(mockDB)

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
	dao := NewUserDAOImpl(mockDB)
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
	dao := NewUserDAOImpl(mockDB)
	timestamp := models.NewDateTimeStamp(now, now)
	user := models.NewUser(-1, "Test User", "Test Email", *timestamp)
	password := "Abcd1234!"
	err := dao.Create(user, password)

	expectedError := "action=UserDAOImpl.Create msg=Error fetching last insert ID: database error"
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
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

	dao := NewUserDAOImpl(mockDB)

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

	dao := NewUserDAOImpl(mockDB)
	user, err := dao.GetById(invalidId)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, fmt.Sprintf("no user found with id %d", invalidId))
	mockDB.AssertExpectations(t)
}
