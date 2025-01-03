package postImpl

import (
	"errors"
	"fmt"
	"go-sns/database"
	"go-sns/models"
	"time"
)



type PostDAOImpl struct{
	db database.Database
}

func NewPostDAOImpl(db database.Database) *PostDAOImpl{
	return &PostDAOImpl{db: db}
}

func (p PostDAOImpl) Create(postData models.Post) error{
	if(postData.GetId() != -1){
		return errors.New("action=PostDAOImpl.Create msg=Cannot create a post data with an existing ID. id: " + string(rune(postData.GetId())))
	}
	
	query := "INSERT INTO posts (title, description, submitted_by, created_at, updated_at) VALUES(?,?,?,?,?)"

	if err := p.db.PrepareAndExecute(query, postData.GetFields()...); err != nil {
		return errors.New("action=PostDAOImpl.Create msg=Error executing query: " + err.Error())
	}

	return nil
}


func (p PostDAOImpl) Delete(id int) error{
	if err := p.db.PrepareAndExecute("DELETE FROM posts WHERE id = ?", id); err != nil{
		return errors.New("action=PostDAOImpl.Delete msg=Error executing query: " + err.Error())	
	}

	return nil
}

func (p PostDAOImpl) GetAll(limitData ...int) ([]models.Post, error){
	var limit int

	recordNums,err := p.db.GetTableLength("posts")
	if err != nil{
		return nil, errors.New("action=PostDAOImpl.GetAll msg=Error executing query: " + err.Error())
	}

	if len(limitData) > 0 && limitData[0] > 0 && limitData[0] <= recordNums{
		limit = limitData[0]
	}else{
		limit = recordNums
	}

	query := "SELECT * FROM posts LIMIT ?"

	posts, err := p.db.PrepareAndFetchAll(query, []interface{}{limit}...)
	if err != nil {
		return nil, errors.New("action=PostDAOImpl.GetAll msg=Error executing query: " + err.Error())
	}

	return p.resultsToPosts(posts), nil
}


func (p PostDAOImpl) GetById(id int) (*models.Post, error){
	results,err := p.db.PrepareAndFetchAll("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return nil, errors.New("action=PostDAOImpl.GetById msg=Error executing query: " + err.Error())
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no post found with id %d", id)
	}

	post := p.resultToPost(results[0])

	return &post, nil
}


func (p PostDAOImpl) resultToPost(post map[string]interface{}) models.Post{
	return *models.NewPost(
		int(post["id"].(int64)),
		int(post["submitted_by"].(int64)),
		post["title"].(string),
		post["description"].(string),
		*models.NewDateTimeStamp(post["created_at"].(time.Time), post["updated_at"].(time.Time)))
}


func (p PostDAOImpl) resultsToPosts(results []map[string]interface{}) []models.Post{
	posts := make([]models.Post, 0)
	
	for _, result := range results{
		posts = append(posts, p.resultToPost(result))
	}

	return posts
}


func (p PostDAOImpl) ValidatePostField(post models.Post) error {
	title := post.GetTitle()
	description := post.GetDescription()

	if len(title) == 0 {
		return errors.New("title is required")
	}
	if len(title) > 100 {
		return errors.New("title must be less than 100 characters")
	}
	if description == "" {
		return errors.New("description is required")
	}
	if len(description) > 1000 {
		return errors.New("description must be less than 1000 characters")
	}
	return nil
}

