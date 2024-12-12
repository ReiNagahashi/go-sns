package implementations

import (
	"errors"
	"go-sns/database"
	"go-sns/models"
	"log"
	"time"
)



type PostDAOImpl struct{}



func (p PostDAOImpl) Create(postData models.Post) bool{
	if(postData.GetId() != -1){
		log.Fatalln("Cannot create a post data with an existing ID. id: " + string(rune(postData.GetId())))
	}
	
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()
	query := "INSERT INTO posts (title, description, created_at, updated_at) VALUES(?,?,?,?)"

	if err := db.PrepareAndExecute(query, postData.GetFields()...); err != nil {
		log.Fatalln("action=PostDAOImpl.Create msg=Error executing query: ", err)
	}

	return true
}


func (p PostDAOImpl) Delete(id int) bool{
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()
	
	if err := db.PrepareAndExecute("DELETE FROM posts WHERE id = ?", id); err != nil{
		log.Fatalln("action=PostDAOImpl.Delete msg=Error executing query: ", err)	
	}

	return true
}

func (p PostDAOImpl) GetAll(limitData ...int) []models.Post{
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var limit int

	if len(limitData) > 1 && limitData[0] > 0{
		limit = limitData[0]
	}else{
		var err error
		limit, err = db.GetTableLength("posts")
		if err != nil{
			log.Fatalln("action=PostDAOImpl.GetAll msg=Error executing query: ", err)
		}
	}

	query := "SELECT * FROM posts LIMIT ?"

	posts, err := db.PrepareAndFetchAll(query, []interface{}{limit}...)
	if err != nil {
		log.Fatalln("action=PostDAOImpl.GetAll msg=Error executing query: ", err)
	}

	return p.resultsToPosts(posts)
}


func (p PostDAOImpl) GetById(id int) []models.Post{
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	post,err := db.PrepareAndFetchAll("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		log.Fatalln("action=PostDAOImpl.GetById msg=Error executing query: ", err)
	}

	return p.resultsToPosts(post)
}


func (p PostDAOImpl) resultToPost(post map[string]interface{}) models.Post{
	return *models.NewPost(
		int(post["id"].(int64)),
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

	if title == "" {
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

