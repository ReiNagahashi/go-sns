package routing

import (
	"encoding/json"
	"fmt"
	"go-sns/config"
	"go-sns/database"
	"go-sns/database/dataAccess/implementations/postImpl"
	"go-sns/database/dataAccess/implementations/userImpl"
	"go-sns/database/dataAccess/interfaces"
	"go-sns/models"
	"go-sns/utils"
	"go-sns/utils/Authenticator"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var id int
	strId := r.URL.Query().Get("id")
	if strId != "" {
		var err error
		id, err = strconv.Atoi(strId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if id != 0 {
		getPostById(w, id)
	} else {
		getAllPosts(w)
	}

}

func getAllPosts(w http.ResponseWriter) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.PostDAO = postImpl.NewPostDAOImpl(db)
	posts, err := dao.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(&posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getPostById(w http.ResponseWriter, id int) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.PostDAO = postImpl.NewPostDAOImpl(db)
	post, err := dao.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.PostDAO = postImpl.NewPostDAOImpl(db)
	timeStamp := time.Now()
	newPost := models.NewPost(-1, r.FormValue("title"), r.FormValue("description"), *models.NewDateTimeStamp(timeStamp, timeStamp))

	if err := dao.ValidatePostField(*newPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := dao.Create(*newPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.PostDAO = postImpl.NewPostDAOImpl(db)

	vars := mux.Vars(r)
	strId := vars["id"]
	if strId == "" {
		APIError(w, "ID is required for retrieve a user", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = dao.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId := vars["id"]
	if strId == "" {
		APIError(w, "ID is required for retrieve a user", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	dao := userImpl.NewUserDAOImpl(db)
	user, err := dao.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	dao := userImpl.NewUserDAOImpl(db)

	timeStamp := time.Now()
	newUser := models.NewUser(-1, r.FormValue("name"), r.FormValue("email"), *models.NewDateTimeStamp(timeStamp, timeStamp))

	if err := dao.ValidateUserField(newUser.GetName(), newUser.GetEmail(), true); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	if err := utils.ValidatePassword(password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := dao.Create(newUser, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := Authenticator.LoginAsUser(newUser, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Login & Logout
func loginHandler(w http.ResponseWriter, r *http.Request) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	dao := userImpl.NewUserDAOImpl(db)

	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := dao.ValidateUserField("", email, false); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := Authenticator.AuthenTicate(email, password, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func isLoggedInHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, err := Authenticator.IsLoggedIn(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(&isLoggedIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := Authenticator.Logout(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func StartWebServer() error {
	r := mux.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	// Post
	r.HandleFunc("/api/posts", postHandler).Methods("GET")                      //全てのポストデータをデータベースから持ってきて表示する。クエリパラメータとしてidが渡されていれば、そのidのデータのみを表示する
	r.HandleFunc("/api/posts", createPostHandler).Methods("POST")               //フォームにタイトル・内容を入力して送信する際に実行されるエンドポイント
	r.HandleFunc("/api/posts/{id:[0-9]+}", deletePostHandler).Methods("DELETE") //各ポストデータに付いている削除ボタンを押したら実行される

	// User
	r.HandleFunc("/api/users/{id:[0-9]+}", userHandler).Methods("GET")
	r.HandleFunc("/api/users/register", registerHandler).Methods("POST")
	// Auth
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/auth/isLoggedIn", isLoggedInHandler).Methods("GET")
	r.HandleFunc("/api/auth/logout", logoutHandler).Methods("POST")

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), corsHandler)
}
