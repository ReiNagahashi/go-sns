package Authenticator

import (
	"errors"
	"go-sns/config"
	"go-sns/database"
	"go-sns/database/dataAccess/implementations/userImpl"
	"go-sns/models"
	"go-sns/utils"
	"log"
	"net/http"
)

var authenticatedUser *models.User

const sessionName = "user_session"

func AuthenTicate(email, password string, w http.ResponseWriter, r *http.Request) (*models.User, error) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	userDao := userImpl.NewUserDAOImpl(db)

	authenticatedUser, err := userDao.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if authenticatedUser == nil {
		return nil, errors.New("could not retrieve user by specified email" + email)
	}
	hashedPassword := userDao.GetHashedPasswordById(authenticatedUser.GetId())

	if utils.CheckPasswordHash(password, hashedPassword) {
		LoginAsUser(authenticatedUser, w, r)

		return authenticatedUser, nil
	}

	return nil, errors.New("invalid password")

}


func LoginAsUser(user *models.User, w http.ResponseWriter, r *http.Request) error {
	if user.GetId() == -1 {
		return errors.New("cannnot login a user with no ID")
	}

	session, err := config.Store.Get(r, sessionName)
	if err != nil {
		return err
	}

	if session.Values["userID"] != nil{
		return errors.New("user is already logged in. Logout before continuing")
	}

	session.Values["userID"] = user.GetId()

	if saveErr := session.Save(r, w); saveErr != nil {
		log.Fatalln(saveErr.Error())
	}

	return nil
}


func Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := config.Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	authenticatedUser = nil
	delete(session.Values, "userID")
	session.Save(r, w)

	return nil
}

func GetAuthenticatedUser(r *http.Request) (*models.User, error) {
	if err := RetrieveAuthenticatedUser(r); err != nil {
		return nil, errors.New("action=RetrieveAuthenticatedUser, Content=" + err.Error())
	}

	return authenticatedUser, nil
}

func RetrieveAuthenticatedUser(r *http.Request) error {
	session, err := config.Store.Get(r, sessionName)
	if err != nil {
		return err
	}

	if session.Values["userID"] == nil {
		return nil
	}

	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	userDao := userImpl.NewUserDAOImpl(db)

	user, err := userDao.GetById(session.Values["userID"].(int))
	if err != nil {
		return err
	}

	authenticatedUser = user

	return nil
}

func IsLoggedIn(r *http.Request) (bool, error) {
	if err := RetrieveAuthenticatedUser(r); err != nil {
		return false, errors.New(err.Error())
	}

	// ユーザーインスタンスが空であるかどうかをチェック
	if authenticatedUser == nil {
		return false, nil
	}

	return true, nil
}
