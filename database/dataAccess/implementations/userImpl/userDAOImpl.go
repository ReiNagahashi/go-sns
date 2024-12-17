package userImpl

import (
	"errors"
	"fmt"
	"go-sns/database"
	"go-sns/models"
	"go-sns/utils"
	"log"
	"time"
)

type UserDAOImpl struct{
	db database.Database
}

func NewUserDAOImpl(db database.Database) *UserDAOImpl{
	return &UserDAOImpl{
		db: db,
	}
}

func (u UserDAOImpl) Create(userData *models.User, password string) error{
	if userData.GetId() != -1{
		return errors.New("action=UserDAOImpl.Create msg=Cannot create a user data with an existing ID. id: " + string(rune(userData.GetId())))
	}

	query := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES(?,?,?,?,?)"

	hashedPassword, err := utils.HashPassword(password)
	if err != nil{
		return errors.New("action=UserDAOImpl.Create msg=Error generating hashpassword: " + err.Error())
	}

	err = u.db.PrepareAndExecute(query, 
			userData.Getname(),
			userData.Getemail(),
			hashedPassword,
			userData.GetTimeStamp().GetCreatedAt(),
			userData.GetTimeStamp().GetUpdatedAt());
	if err != nil {
		return errors.New("action=UserDAOImpl.Create msg=Error executing query: " + err.Error())
	}

	lastInsertId, err := u.db.GetLastInsertedId()
	if err != nil{
		return errors.New("action=UserDAOImpl.Create msg=Error fetching last insert ID: " + err.Error())
	}

	userData.SetId(lastInsertId)

	return nil
}


func (u UserDAOImpl) getRawById(id int) map[string]interface{}{
	query := "SELECT * FROM users WHERE id=?"
	result,err := u.db.PrepareAndFetchAll(query, id)

	if err != nil{
		log.Fatalln("action=UserDAOImpl.getRawById msg=Error executing query: ", err)
	}

	return result[0]
}

func (u UserDAOImpl) GetByEmail(email string) (*models.User, error){
	result,err := u.db.PrepareAndFetchAll("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, errors.New("action=UserDAOImpl.GetByEmail msg=Error executing query: "+err.Error())
	}
	user := u.resultToUser(result[0])

	return &user, nil
}

func (u UserDAOImpl) GetById(id int) (*models.User, error){
	results,err := u.db.PrepareAndFetchAll("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, errors.New("action=UserDAOImpl.GetById msg=Error executing query: "+err.Error())
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no user found with id %d", id)
	}

	user := u.resultToUser(results[0])

	return &user, nil
}


func (u UserDAOImpl) GetHashedPasswordById(id int)string{
	v := u.getRawById(id)["password"]

	return v.(string)
}


func (u UserDAOImpl) resultToUser(user map[string]interface{}) models.User{
	return *models.NewUser(
		int(user["id"].(int)),
		user["name"].(string),
		user["email"].(string),
		*models.NewDateTimeStamp(user["created_at"].(time.Time), user["updated_at"].(time.Time)))
}


// func (u UserDAOImpl) resultsToUsers(results []map[string]interface{}) []models.User{
// 	users := make([]models.User, 0)
	
// 	for _, result := range results{
// 		users = append(users, u.resultToUser(result))
// 	}

// 	return users
// }


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

