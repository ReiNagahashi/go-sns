package routing

import (
	"encoding/json"
	"fmt"
	"go-sns/config"
	"go-sns/database/dataAccess/implementations"
	"log"
	"net/http"
	"regexp"
)

type JSONError struct {
	Error string `json:"error"`
	Code int `json:"code"`
}

func APIError(w http.ResponseWriter, errMessage string, code int){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil{
		log.Fatal(err)
	}
	w.Write(jsonError)
}


var apiValidaPath = regexp.MustCompile("^/api/posts/$")

// URLのバリデーション
func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		m := apiValidaPath.FindStringSubmatch(r.URL.Path)
		if len(m) == 0{
			APIError(w, "Not found", http.StatusNotFound)
		}
		fn(w, r)
	}
}


func viewPostsHandler(w http.ResponseWriter, r *http.Request){
	dao := implementations.PostDAOImpl{}
	posts := dao.GetAll()

	fmt.Println(posts)

	js, err := json.Marshal(posts)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func StartWebServer() error {
	http.HandleFunc("/api/posts/", apiMakeHandler(viewPostsHandler))

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}