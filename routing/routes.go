package routing

import (
	"encoding/json"
	"fmt"
	"go-sns/config"
	"log"
	"net/http"
	"regexp"
)

type JSONError struct {
	Error string `json:"error"`
	Code int `json:"code"`
}

// api通信をしたときのエラーを表示するためにJSON形式のエラーを作成する
func APIError(w http.ResponseWriter, errMessage string, code int){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil{
		log.Fatal(err)
	}
	w.Write(jsonError)
}


var apiValidaPath = regexp.MustCompile("^/api/candle/$")

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
	// ポストインスタンスを全て持ってくる
	// エンコード
	// ヘッダーの設定
	// ヘッダーにレスポンスデータの追加
}


func StartWebServer() error {
	http.HandleFunc("/api/posts/", apiMakeHandler(viewPostsHandler))

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}