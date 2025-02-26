package utils

import (
	"errors"
	"go-sns/config"
	"net/http"
)

func GetCsrfToken(w http.ResponseWriter, r *http.Request) (string, error){
	session, err := config.Store.Get(r, "csrf_token")
	if err != nil{
		return "", err
	}
	
	if session.Values["token"] == nil{
		return "", errors.New("csrf_token is not generated")
	}

	return session.Values["token"].(string), nil
}


