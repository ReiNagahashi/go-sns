package main

import (
	"go-sns/config"
	"go-sns/routing"
	"go-sns/utils"
)

func main(){
	utils.LoggingSettings(config.Config.LogFile)
	routing.StartWebServer()
}
