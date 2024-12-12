package routing

import (
	"encoding/json"
	"fmt"
	"go-sns/config"
	"go-sns/database/dataAccess/implementations"
	"go-sns/database/dataAccess/interfaces"
	"go-sns/models"
	"go-sns/utils"
	"go-sns/utils/Authenticator"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

func postHandler(w http.ResponseWriter, r *http.Request){
	var id int
	strId := r.URL.Query().Get("id")
	if strId != ""{
		var err error
		id,err = strconv.Atoi(strId)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if id != 0{
		getPostById(w, id)
	}else{
		getAllPosts(w)
	}

}


func getAllPosts(w http.ResponseWriter){
	var dao interfaces.PostDAO = implementations.PostDAOImpl{}
	posts := dao.GetAll()

	js, err := json.Marshal(&posts)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func getPostById(w http.ResponseWriter, id int){
	var dao interfaces.PostDAO = implementations.PostDAOImpl{}
	post := dao.GetById(id)

	js, err := json.Marshal(&post)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func createPostHandler (w http.ResponseWriter, r *http.Request){
	var dao interfaces.PostDAO = implementations.PostDAOImpl{}

	timeStamp := time.Now()
	newPost := models.NewPost(-1, r.FormValue("title"), r.FormValue("description"), *models.NewDateTimeStamp(timeStamp, timeStamp))

	if err := dao.ValidatePostField(*newPost); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	success := dao.Create(*newPost)
	if !success{
		APIError(w, "Post creation failed", http.StatusInternalServerError)
		return
	}

}


func deletePostHandler (w http.ResponseWriter, r *http.Request){
	var dao interfaces.PostDAO = implementations.PostDAOImpl{}

	vars := mux.Vars(r)
	strId := vars["id"]
	if strId == ""{
		APIError(w, "ID is required for retrieve a user", http.StatusBadRequest)
		return
	}

	id,err := strconv.Atoi(strId)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	success := dao.Delete(id)
	if !success{
		APIError(w, "Data deletion failed", http.StatusInternalServerError)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	strId := vars["id"]
	if strId == ""{
		APIError(w, "ID is required for retrieve a user", http.StatusBadRequest)
		return
	}

	id,err := strconv.Atoi(strId)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var dao interfaces.UserDAO = implementations.UserDAOImpl{}
	user := dao.GetById(id)

	js, err := json.Marshal(&user)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}


func registerHandler(w http.ResponseWriter, r *http.Request){
	var dao interfaces.UserDAO = implementations.UserDAOImpl{}

	timeStamp := time.Now()
	newUser := models.NewUser(-1, r.FormValue("name"), r.FormValue("email"), *models.NewDateTimeStamp(timeStamp, timeStamp))

	if err := dao.ValidateUserField(newUser.Getname(), newUser.Getemail()); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	if err := utils.ValidatePassword(password); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	success := dao.Create(newUser, password)
	if !success{
		APIError(w, "User creation failed", http.StatusInternalServerError)
		return
	}

	if err := Authenticator.LoginAsUser(*newUser, w, r); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


// Login & Logout
func loginHandler(w http.ResponseWriter, r *http.Request){
	userDao := implementations.UserDAOImpl{}
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := userDao.ValidateUserField(name, email); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	_, err := Authenticator.AuthenTicate(email, password, w, r)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request){
	if err := Authenticator.Logout(w, r); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("User Logout successfully!!!")
}


func StartWebServer() error {
	r := mux.NewRouter()
	// Post
	r.HandleFunc("/api/posts", postHandler).Methods("GET")
	r.HandleFunc("/api/posts", createPostHandler).Methods("POST")
	r.HandleFunc("/api/posts/{id:[0-9]+}", deletePostHandler).Methods("DELETE")

	// User
	r.HandleFunc("/api/users/{id:[0-9]+}", userHandler).Methods("GET")
	r.HandleFunc("/api/users/register", registerHandler).Methods("POST")
	// Auth
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/auth/logout", logoutHandler).Methods("POST")

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), r)
}