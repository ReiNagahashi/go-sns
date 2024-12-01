package main

import (
	"go-sns/config"
	"go-sns/utils"
	"log"
)

func main(){
	utils.LoggingSettings(config.Config.LogFile)
	log.Println("test")
}