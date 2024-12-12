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

func (u UserDAOImpl) Create(userData *models.User, password string) bool{
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

	var lastInsertId int64
	err = db.DbConnection.QueryRow("SELECT last_insert_rowid()").Scan(&lastInsertId)
	if err != nil{
		log.Fatalln("action=UserDAOImpl.Create msg=Error retrieving last insert ID: ", err)
	}

	userData.SetId(int(lastInsertId))

	return true
}

func (u UserDAOImpl) getRawById(id int) map[string]interface{}{
	db := database.NewSqliteBase()
	query := "SELECT * FROM users WHERE id=?"
	result,err := db.PrepareAndFetchAll(query, id)

	if err != nil{
		log.Fatalln("action=UserDAOImpl.getRawById msg=Error executing query: ", err)
	}

	return result[0]
}

func (u UserDAOImpl) GetByEmail(email string) models.User{
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	user,err := db.PrepareAndFetchAll("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		log.Fatalln("action=UserDAOImpl.GetByEmail msg=Error executing query: ", err)
	}

	return u.resultsToUsers(user)[0]
}

func (u UserDAOImpl) GetById(id int) models.User{
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	user,err := db.PrepareAndFetchAll("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatalln("action=UserDAOImpl.GetById msg=Error executing query: ", err)
	}

	return u.resultsToUsers(user)[0]
}


func (u UserDAOImpl) GetHashedPasswordById(id int)string{
	v := u.getRawById(id)["password"]

	return v.(string)
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


func (u UserDAOImpl) ValidateUserField(fields ...interface{}) error {
	name := fields[0].(string)
	email := fields[1].(string)

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

