package routing

import (
	"encoding/json"
	"fmt"
	"go-sns/config"
	"go-sns/database"
	"go-sns/database/dataAccess/implementations/postImpl"
	"go-sns/database/dataAccess/implementations/userImpl"
	"go-sns/database/dataAccess/interfaces"
	"go-sns/middleware"
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

func addFavoritePostHandler(w http.ResponseWriter, r *http.Request) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var postDao interfaces.PostDAO = postImpl.NewPostDAOImpl(db)

	vars := mux.Vars(r)
	strId := vars["id"]
	if strId == "" {
		APIError(w, "ID is required to add a favorite", http.StatusBadRequest)
		return
	}

	postId, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ログインしているユーザーを取得
	authUser, err := Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = postDao.AddFavorite(authUser.GetId(), postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.PostDAO = postImpl.NewPostDAOImpl(db)

	vars := mux.Vars(r)
	strId := vars["id"]
	if strId == "" {
		APIError(w, "Post id is required to add a favorite", http.StatusBadRequest)
		return
	}

	postId, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ログインしているユーザーを取得
	authUser, err := Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if authUser == nil {
		APIError(w, "User is invalid", http.StatusBadRequest)
		return
	}

	err = dao.AddFavorite(authUser.GetId(), postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.PostDAO = postImpl.NewPostDAOImpl(db)

	vars := mux.Vars(r)
	strId := vars["id"]
	if strId == "" {
		APIError(w, "Post id is required to remove a favorite", http.StatusBadRequest)
		return
	}

	postId, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ログインしているユーザーを取得
	authUser, err := Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if authUser == nil {
		APIError(w, "User is invalid", http.StatusBadRequest)
		return
	}

	err = dao.DeleteFavorite(authUser.GetId(), postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var id, limit int
	strId := r.URL.Query().Get("id")
	strLimit := r.URL.Query().Get("limit")

	if strId != "" {
		var err error
		id, err = strconv.Atoi(strId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if strLimit != "" {
		var err error
		limit, err = strconv.Atoi(strLimit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.PostDAO = postImpl.NewPostDAOImpl(db)

	if id != 0 && limit == 0 {
		getPostById(w, id, dao)
	} else if limit != 0 {
		getPostWithLimit(w, limit, dao)
	} else {
		getAllPosts(w, dao)
	}

}

func getAllPosts(w http.ResponseWriter, dao interfaces.PostDAO) {
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

func getPostWithLimit(w http.ResponseWriter, limit int, dao interfaces.PostDAO) {
	posts, err := dao.GetPosts(limit)
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

func getPostById(w http.ResponseWriter, id int, dao interfaces.PostDAO) {
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
	authUser, err := Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newPost := models.NewPost(-1, authUser.GetId(), r.FormValue("title"), r.FormValue("description"), *models.NewDateTimeStamp(timeStamp, timeStamp), []models.User{})

	if err := dao.ValidatePostField(*newPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = dao.Create(*newPost)
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
		APIError(w, "ID is required to delete a post", http.StatusBadRequest)
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
		getUserById(w, id)
	} else {
		getUsers(w)
	}
}

func getUsers(w http.ResponseWriter) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.UserDAO = userImpl.NewUserDAOImpl(db)
	users, err := dao.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getUserById(w http.ResponseWriter, id int) {
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()

	var dao interfaces.UserDAO = userImpl.NewUserDAOImpl(db)
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

	var dao interfaces.UserDAO = userImpl.NewUserDAOImpl(db)

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

	var dao interfaces.UserDAO = userImpl.NewUserDAOImpl(db)

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

func getLoggedinUserHandler(w http.ResponseWriter, r *http.Request) {
	loggedinUser, err := Authenticator.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(loggedinUser)
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

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Post
	r.HandleFunc("/api/posts", postHandler).Methods("GET")                                   //全てのポストデータをデータベースから持ってきて表示する。クエリパラメータとしてidが渡されていれば、そのidのデータのみを表示する
	r.HandleFunc("/api/posts", createPostHandler).Methods("POST")                            //フォームにタイトル・内容を入力して送信する際に実行されるエンドポイント
	r.HandleFunc("/api/posts/favorite/{id:[0-9]+}", addFavoritePostHandler).Methods("POST")  //各ポストデータに付いているハートアイコンボタンを押したら実行される
	r.HandleFunc("/api/posts/{id:[0-9]+}", deletePostHandler).Methods("DELETE")              //各ポストデータに付いている削除ボタンを押したら実行される
	r.HandleFunc("/api/posts/favorite/{id:[0-9]+}", addFavoriteHandler).Methods("POST")      //各ポストデータに付いている削除ボタンを押したら実行される
	r.HandleFunc("/api/posts/favorite/{id:[0-9]+}", deleteFavoriteHandler).Methods("DELETE") //各ポストデータに付いている削除ボタンを押したら実行される

	// User
	r.HandleFunc("/api/users", userHandler).Methods("GET")
	r.HandleFunc("/api/users/register", registerHandler).Methods("POST")
	// Auth
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/auth/loggedInUser", getLoggedinUserHandler).Methods("GET")
	r.HandleFunc("/api/auth/logout", logoutHandler).Methods("POST")

	wrappedMux := middleware.NewCsrfMiddleware(corsHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), wrappedMux)
}
