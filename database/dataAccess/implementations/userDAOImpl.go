package implementations

import (
	"errors"
	"go-sns/database"
	"go-sns/models"
	"go-sns/utils"
	"log"
	"time"
)



type UserDAOImpl struct{}

func (u UserDAOImpl) Create(userData models.User, password string) bool{
	if userData.GetId() != -1{
		log.Fatalln("Cannot create a user data with an existing ID. id: " + string(rune(userData.GetId())))
	}

	db := database.NewSqliteBase()
	defer db.DbConnection.Close()
	query := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES(?,?,?,?,?)"

	hashedPassword, err := utils.HashPassword(password)
	if err != nil{
		log.Fatalln("action=UserDAOImpl.Create msg=Error executing query: ", err)
	}

	err = db.PrepareAndExecute(query, 
			userData.Getname(),
			userData.Getemail(),
			hashedPassword,
			userData.GetTimeStamp().GetCreatedAt(),
			userData.GetTimeStamp().GetUpdatedAt());
	if err != nil {
		log.Fatalln("action=UserDAOImpl.Create msg=Error executing query: ", err)
	}

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


func (u UserDAOImpl) ValidateUserField(user models.User) error {
	name := user.Getname()
	email := user.Getemail()

	if name == "" {
		return errors.New("name is required")
	}
	if len(name) > 20 {
		return errors.New("name must be less than 20 characters")
	}
	if email == "" {
		return errors.New("email is required")
	}
	if len(email) > 30 {
		return errors.New("email must be less than 30 characters")
	}

	return nil
}

