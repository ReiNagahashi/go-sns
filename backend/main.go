package main

import (
	"go-sns/config"
	"go-sns/routing"
	"go-sns/utils"
	"log"
)

func main(){
	utils.LoggingSettings(config.Config.LogFile)
	utils.RunMigrations()
	utils.RunSeedings()
	log.Println(routing.StartWebServer())
}
