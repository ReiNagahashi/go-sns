package main

import (
	"go-sns/config"
	"go-sns/routing"
	"go-sns/utils"
	"log"
)

// TODO Postマンで登録、ログインをフォームデータでテストする
func main(){
	utils.LoggingSettings(config.Config.LogFile)
	log.Println(routing.StartWebServer())
}
