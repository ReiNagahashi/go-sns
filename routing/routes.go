package routing

import (
	"encoding/json"
	"fmt"
	"go-sns/config"
	"go-sns/database/dataAccess/implementations"
	"go-sns/database/dataAccess/interfaces"
	"go-sns/models"
	"go-sns/utils"
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
			if id != 0{
				getPostById(w, id)
			}else{
				getAllPosts(w)
			}
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
	var dao interfaces.PostDAO = implementations.PostDAOImpl{}
	posts := dao.GetAll()

	js, err := json.Marshal(posts)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func getPostById(w http.ResponseWriter, id int){
	var dao interfaces.PostDAO = implementations.PostDAOImpl{}
	post := dao.GetById(id)

	js, err := json.Marshal(post)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func createPost (w http.ResponseWriter, r *http.Request){
	var dao interfaces.PostDAO = implementations.PostDAOImpl{}

	timeStamp := time.Now()
	newPost := models.NewPost(-1, r.URL.Query().Get("title"), r.URL.Query().Get("description"), *models.NewDateTimeStamp(timeStamp, timeStamp))

	if err := dao.ValidatePostField(*newPost); err != nil{
		APIError(w, "Validation Error: " + err.Error(), http.StatusBadRequest)
	}
	success := dao.Create(*newPost)
	if !success{
		APIError(w, "Post creation failed", http.StatusInternalServerError)
	}

}


func deletePost (w http.ResponseWriter, id int){
	var dao interfaces.PostDAO = implementations.PostDAOImpl{}

		success := dao.Delete(id)
		if !success{
			APIError(w, "Data deletion failed", http.StatusInternalServerError)
		}
}


var apiUserPath = regexp.MustCompile("^/api/users(/([0-9]+))?/?$")

func userHandler(w http.ResponseWriter, r *http.Request){
	matches := apiUserPath.FindStringSubmatch(r.URL.Path)
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
			if id == 0{
				APIError(w, "ID is required for retrieve a user", http.StatusBadRequest)
				return
			}
			getUserById(w, id)
		case http.MethodPost:
			createUser(w, r)
		default:
			APIError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}


func createUser(w http.ResponseWriter, r *http.Request){
	var dao interfaces.UserDAO = implementations.UserDAOImpl{}

	timeStamp := time.Now()
	newUser := models.NewUser(-1, r.URL.Query().Get("name"), r.URL.Query().Get("email"), *models.NewDateTimeStamp(timeStamp, timeStamp))

	if err := dao.ValidateUserField(*newUser); err != nil{
		APIError(w, "Validation Error: " + err.Error(), http.StatusBadRequest)
	}

	password := r.URL.Query().Get("password")
	if err := utils.ValidatePassword(password); err != nil{
		APIError(w, "Validation Error: " + err.Error(), http.StatusBadRequest)
	}

	success := dao.Create(*newUser, password)
	if !success{
		APIError(w, "User creation failed", http.StatusInternalServerError)
	}
}

func getUserById(w http.ResponseWriter, id int){
	var dao interfaces.UserDAO = implementations.UserDAOImpl{}
	user := dao.GetById(id)

	js, err := json.Marshal(user)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func StartWebServer() error {
	http.HandleFunc("/api/posts", postHandler)
	http.HandleFunc("/api/users", userHandler)


	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}