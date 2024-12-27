package config

import (
	"log"
	"os"
	"github.com/gorilla/sessions"
	"gopkg.in/ini.v1"
)

type ConfigList struct{
	LogFile string
	SQLDriver string
	DbPath string
	Port int
	Session_key string
}

var Config ConfigList

var Store *sessions.CookieStore


func init(){
    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        configPath = "config.ini" 
    }

    cfg, err := ini.Load(configPath)
    if err != nil {
        log.Printf("Failed to read file: %v", err)
        os.Exit(1)
    }

	Config = ConfigList{
		LogFile: cfg.Section("go-sns").Key("log_file").String(),
		SQLDriver: cfg.Section("db").Key("driver").String(),
		DbPath: cfg.Section("db").Key("path").String(),
		Port: cfg.Section("web").Key("port").MustInt(),
		Session_key: cfg.Section("web").Key("session_key").String(),
	}

	Store = sessions.NewCookieStore([]byte(Config.Session_key))

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,  // クッキーをJavaScriptからアクセス不能にする
		Secure:   false, // HTTPSのみでセッションクッキーを送信する
	}
}