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
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")

    // 🔹 OPTIONS メソッドの場合は 200 OK を返す
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
	session, err := config.Store.Get(r, "csrf_token")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var token string

	if session.Values["token"] == nil{
		bytes := make([]byte, 32)
		_, err := rand.Read(bytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token = hex.EncodeToString(bytes)

		session.Values["token"] = token
	}else{
		token = session.Values["token"].(string)
	}

	http.SetCookie(w, &http.Cookie{
		Name: "csrf_token",
		Value: token,
		HttpOnly: false,
		Secure: false,
		Path: "/",
		SameSite: http.SameSiteLaxMode, // CSRF対策のため適切に設定
	})

	err = session.Save(r, w)
	if err != nil{
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}


	if r.Method != "GET" && token != r.Header.Get("X-CSRF-Token"){
		http.Error(w,  "CSRF token is invalid. Access has been denied.", http.StatusBadRequest)
		return
	}    


	l.handler.ServeHTTP(w, r)
}


func NewCsrfMiddleware(handlerToWrap http.Handler) *csrfMiddleware {
	return &csrfMiddleware{handlerToWrap}
}
