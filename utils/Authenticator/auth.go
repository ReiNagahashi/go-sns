package Authenticator

import (
	"errors"
	"go-sns/config"
	"go-sns/database/dataAccess/implementations"
	"go-sns/models"
	"go-sns/utils"
	"net/http"
	"unsafe"

	"github.com/gorilla/sessions"
)


var authenticatedUser models.User
var store = sessions.NewCookieStore([]byte(config.Config.Session_key))
const sessionName = "user-session"


func AuthenTicate(email, password string, w http.ResponseWriter, r *http.Request) (*models.User, error){
	userDao := implementations.UserDAOImpl{}
	authenticatedUser = userDao.GetByEmail(email)

	if unsafe.Sizeof(authenticatedUser) == 0{
		return nil, errors.New("could not retrieve user by specified email" + email)
	}
	hashedPassword := userDao.GetHashedPasswordById(authenticatedUser.GetId())

	if utils.CheckPasswordHash(password, hashedPassword){
		LoginAsUser(authenticatedUser, w, r)

		return &authenticatedUser, nil
	}

	return nil, errors.New("invalid password")
	
}

func LoginAsUser(user models.User, w http.ResponseWriter, r *http.Request) error{
	if user.GetId() == -1{
		return errors.New("cannnot login a user with no ID")
	}
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values["userID"] = user.GetId()
	session.Save(r, w)

	return nil
}


func Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, sessionName)
	if err != nil{
		return err
	}
	delete(session.Values, "userID")
	session.Save(r, w)

	return nil
}


func RetrieveAuthenticatedUser(r *http.Request)error {
	session, err := store.Get(r, sessionName)
	if err != nil{
		return err
	}

	if session.Values["userID"] == nil{
		return nil
	}
	userDao := implementations.UserDAOImpl{}

	authenticatedUser = userDao.GetById(session.Values["userID"].(int))

	return nil
}


func IsLoggedIn(r *http.Request) (bool, error){
	if err := RetrieveAuthenticatedUser(r); err != nil{
		return false, errors.New(err.Error())
	}
	// ユーザーインスタンスが空であるかどうかをチェック
	if unsafe.Sizeof(authenticatedUser) == 0{
		return false, nil
	}

	return true, nil
}