package postImpl

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

// カラム内容に関係なく、PrepareAndExecuteがdao.Createによって呼び出されたかどうかをチェックする
func TestCreate_Success(t *testing.T){
	mockDB := new(mocks.Database)
	// mock.Anythingとすることで、カラム内容を重視ししないことにしている
    mockDB.On(
        "PrepareAndExecute",
        mock.Anything,
        mock.Anything,
        mock.Anything,
        mock.Anything,
        mock.Anything,
        mock.Anything,
    ).Return(nil)
	
	now := time.Now().Truncate(time.Second)
	timestamp := models.NewDateTimeStamp(now, now)
	post := models.NewPost(-1, 1, "Test Title", "Test Description", *timestamp, []models.User{})
	
	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)
	err := dao.Create(*post)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}


func TestCreate_InvalidID(t *testing.T){
	mockDB := new(mocks.Database)
	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)

	invalidID := 1
	post := models.NewPost(invalidID, 1, "Test Title", "Test Description", models.DateTimeStamp{}, []models.User{})

	err := dao.Create(*post)

	expectedError := "action=PostDAOImpl.Create msg=Cannot create a post data with an existing ID. id: "  + string(rune(invalidID))
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	mockDB.AssertNotCalled(t, "PrepareAndExecute")
}


func TestDelete_Success(t *testing.T){
	mockDB := new(mocks.Database)
	mockDB.On(
		"PrepareAndExecute",
		mock.Anything,
		mock.Anything).
			Return(nil)
	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)

	id := 1
	err := dao.Delete(id)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}


func TestDelete_QueryExecutionFailure(t *testing.T) {
	mockDB := new(mocks.Database)

	queryError := errors.New("database error")
	mockDB.On(
		"PrepareAndExecute",
		"DELETE FROM posts WHERE id = ?",
		-10).
			Return(queryError)

	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)

	id := -10
	err := dao.Delete(id)

	expectedError := "action=PostDAOImpl.Delete msg=Error executing query: database error"
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	mockDB.AssertExpectations(t)
}


func generateMockPosts(count int, now time.Time) []map[string]interface{} {
	mockData := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		mockData[i] = map[string]interface{}{
			"id":          i + 1,
			"title":       fmt.Sprintf("Test Title %d", i), 
			"description": fmt.Sprintf("Test Description %d", i),
			"submitted_by": 1,
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
		"posts",).Return(recordSize, nil)

	mockDB.On(
		"PrepareAndFetchAll",
		mock.Anything,
		recordSize).Return(generateMockPosts(recordSize, now), nil)

	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)

	posts, err := dao.GetAll()

	assert.NoError(t, err)
	assert.Len(t, posts, recordSize)
	titleId := 10
	descriptionId := 50
	expectedTitle := fmt.Sprintf("Test Title %v", titleId)
	expecteDescription := fmt.Sprintf("Test Description %v",descriptionId)

	assert.Equal(t, expectedTitle, posts[titleId].GetTitle())
	assert.Equal(t, expecteDescription,posts[descriptionId].GetDescription())

	mockDB.AssertExpectations(t)
}

func TestGetAllWithLimit_Success(t *testing.T){
	var limit int = 10
	var recordSize int = 100
	
    now := time.Now().Truncate(time.Second)

	mockDB := new(mocks.Database)

	mockDB.On(
		"GetTableLength",
		"posts",).Return(recordSize, nil)

	mockDB.On(
		"PrepareAndFetchAll",
		mock.Anything,
		limit).Return(generateMockPosts(limit,now), nil)

	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)

	posts, err := dao.GetAll(limit)

	assert.NoError(t, err)
	assert.Len(t, posts, limit)

	titleId := 5
	descriptionId := 3
	expectedTitle := fmt.Sprintf("Test Title %v", titleId)
	expecteDescription := fmt.Sprintf("Test Description %v",descriptionId)

	assert.Equal(t, expectedTitle, posts[titleId].GetTitle())
	assert.Equal(t, expecteDescription,posts[descriptionId].GetDescription())

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
		"posts",).Return(recordSize, nil)

	mockDB.On(
		"PrepareAndFetchAll",
		mock.Anything,
		recordSize).Return(generateMockPosts(recordSize,now), nil)

	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)
	posts, err := dao.GetAll(limit)

	assert.NoError(t, err)
	assert.Len(t, posts, recordSize)
	titleId := 10
	descriptionId := 50
	expectedTitle := fmt.Sprintf("Test Title %v", titleId)
	expecteDescription := fmt.Sprintf("Test Description %v",descriptionId)

	assert.Equal(t, expectedTitle, posts[titleId].GetTitle())
	assert.Equal(t, expecteDescription,posts[descriptionId].GetDescription())

	mockDB.AssertExpectations(t)
}


func TestGetById_Success(t *testing.T){
	var id int = 10
	expectedTitle := fmt.Sprintf("Test Title %v", id)
    now := time.Now().Truncate(time.Second)
	mockDB := new(mocks.Database)
	mockDB.On(
		"PrepareAndFetchAll",
		mock.Anything,
		id).Return([]map[string]interface{}{
			{
				"id":          id,
				"title":       expectedTitle,
				"submitted_by": 1,
				"description": "",
				"created_at":  now,
				"updated_at":  now,
			},
		}, nil)

	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)

	post, err := dao.GetById(id)

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, id, post.GetId())
	assert.Equal(t, expectedTitle, post.GetTitle())

	mockDB.AssertExpectations(t)
}


func TestGetById_NotFound(t *testing.T){
	var invalidId int = 100000

	mockDB := new(mocks.Database)
	mockDB.On("PrepareAndFetchAll", mock.Anything, invalidId).
		Return([]map[string]interface{}{}, nil)
	
	var dao interfaces.PostDAO = NewPostDAOImpl(mockDB)
	post, err := dao.GetById(invalidId)

	assert.Error(t, err)
	assert.Nil(t, post)
	assert.EqualError(t, err, fmt.Sprintf("no post found with id %d", invalidId))
	mockDB.AssertExpectations(t)
}