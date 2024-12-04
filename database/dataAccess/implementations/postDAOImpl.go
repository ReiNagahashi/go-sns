package implementations

import (
	"go-sns/database"
	"go-sns/models"
	"log"
)



type PostDAOImpl struct{

}

func (p PostDAOImpl) Create(postData models.Post) bool{
	if(postData.GetId() != -1){
		log.Fatalln("Cannot create a computer part with an existing ID. id: " + string(postData.GetId()))
	}

	db := database.NewSqliteBase()

	defer db.DbConnection.Close()
	query := "INSERT INTO posts (title, description, created_at, updated_at) VALUES(?,?,?,?)"

	if err := db.Execute(query, postData.GetFields()...); err != nil {
		log.Fatalln("action=PostDAOImpl.Create msg=Error executing query: ", err)
	}

	return true

}
