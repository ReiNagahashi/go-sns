package routing

import (
	"encoding/json"
	"fmt"
	"go-sns/config"
	"go-sns/database/dataAccess/implementations"
	"go-sns/models"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
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


var apiPostPath = regexp.MustCompile("^/api/posts(/([0-9]+))?/?$")

func postHandler(w http.ResponseWriter, r *http.Request){
	matches := apiPostPath.FindStringSubmatch(r.URL.Path)
	if matches == nil{
		APIError(w, "Not found", http.StatusNotFound)
		return
	}

	var id int
	strId := r.URL.Query().Get("id")
	if strId != ""{
		var err error
		id,err = strconv.Atoi(strId)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	switch r.Method {
		case http.MethodGet:
			getAllPosts(w)
		case http.MethodPost:
			createPost(w, r)
		case http.MethodDelete:
			if id == 0{
				APIError(w, "ID is required for delete", http.StatusBadRequest)
				return
			}
			deletePost(w, id)
		default:
			APIError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}


func getAllPosts(w http.ResponseWriter){
	dao := implementations.PostDAOImpl{}
	posts := dao.GetAll()

	js, err := json.Marshal(posts)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func createPost (w http.ResponseWriter, r *http.Request){
	dao := implementations.PostDAOImpl{}

	timeStamp := time.Now()
	newPost := models.NewPost(-1, r.URL.Query().Get("title"), r.URL.Query().Get("description"), *models.NewDateTimeStamp(timeStamp, timeStamp))

	if err := dao.ValidatePost(*newPost); err != nil{
		APIError(w, "Validation Error: " + err.Error(), http.StatusBadRequest)
	}
	success := dao.Create(*newPost)
	if !success{
		APIError(w, "Data deletion failed", http.StatusInternalServerError)
	}

}


func deletePost (w http.ResponseWriter, id int){
		dao := implementations.PostDAOImpl{}

		success := dao.Delete(id)
		if !success{
			APIError(w, "Data deletion failed", http.StatusInternalServerError)
		}
}


func StartWebServer() error {
	http.HandleFunc("/api/posts", postHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}