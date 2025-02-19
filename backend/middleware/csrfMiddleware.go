package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"go-sns/config"
	"net/http"
)

type csrfMiddleware struct {
	handler http.Handler
}

func (l *csrfMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, "csrf_token")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["token"] == nil{
		bytes := make([]byte, 32)
		_, err := rand.Read(bytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["token"] = hex.EncodeToString(bytes)
		session.Save(r, w)
	}

	correctToken := session.Values["token"]

	if r.Method != "GET"{
		if correctToken != r.FormValue("csrf_token"){
			http.Error(w,  "CSRF token is invalid. Access has been denied.", http.StatusBadRequest)
			return
		}
	}

	l.handler.ServeHTTP(w, r)
}


func NewCsrfMiddleware(handlerToWrap http.Handler) *csrfMiddleware {
	return &csrfMiddleware{handlerToWrap}
}
