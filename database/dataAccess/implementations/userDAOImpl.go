package implementations

import (
	"go-sns/database"
	"go-sns/models"
	"log"
	"time"
)



type UserDAOImpl struct{}

func (u UserDAOImpl) Create(userData models.User) bool{
	return true
}
func (u UserDAOImpl) GetByEmail(email string) models.User{
	return models.User{}
}
func (u UserDAOImpl) GetHashedPasswordById(id int) string{
	return ""
}

func (u UserDAOImpl) GetById(id int) []models.User{
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	user,err := db.PrepareAndFetchAll("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatalln("action=UserDAOImpl.GetById msg=Error executing query: ", err)
	}

	return u.resultsToUsers(user)
}

func (u UserDAOImpl) resultToUser(user map[string]interface{}) models.User{
	return *models.NewUser(
		int(user["id"].(int64)),
		user["name"].(string),
		user["email"].(string),
		*models.NewDateTimeStamp(user["created_at"].(time.Time), user["updated_at"].(time.Time)))
}


func (u UserDAOImpl) resultsToUsers(results []map[string]interface{}) []models.User{
	users := make([]models.User, 0)
	
	for _, result := range results{
		users = append(users, u.resultToUser(result))
	}

	return users
}