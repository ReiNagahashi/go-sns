package implementations

import (
	"go-sns/database"
	"go-sns/models"
	"log"
	"time"
)



type PostDAOImpl struct{

}


func (p PostDAOImpl) GetAll(offset, limit int) []models.Post{
	db := database.NewSqliteBase()

	defer db.DbConnection.Close()
	query := "SELECT * FROM posts LIMIT ?, ?"

	posts, err := db.PrepareAndFetchAll(query, []interface{}{offset, limit}...)
	if err != nil {
		log.Fatalln("action=PostDAOImpl.GetAll msg=Error executing query: ", err)
	}

	return p.resultsToPosts(posts)
}


func (p PostDAOImpl) Create(postData models.Post) bool{
	if(postData.GetId() != -1){
		log.Fatalln("Cannot create a computer part with an existing ID. id: " + string(postData.GetId()))
	}

	db := database.NewSqliteBase()

	defer db.DbConnection.Close()
	query := "INSERT INTO posts (title, description, created_at, updated_at) VALUES(?,?,?,?)"

	if err := db.PrepareAndExecute(query, postData.GetFields()...); err != nil {
		log.Fatalln("action=PostDAOImpl.Create msg=Error executing query: ", err)
	}

	return true
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

// func (p PostDAOImpl) Delete(postData models.Post) bool{

// }
