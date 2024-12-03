package main

import (
	"fmt"
	"go-sns/config"
	"go-sns/database"
	"go-sns/utils"
)

func main(){
	utils.LoggingSettings(config.Config.LogFile)
	// log.Println("test")
	fmt.Println(database.DbConnection)
}