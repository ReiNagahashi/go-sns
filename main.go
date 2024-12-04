package main

import (
	"fmt"
	"go-sns/config"
	"go-sns/database/dataAccess/implementations"
	"go-sns/models"
	"go-sns/utils"
	"time"
)

func main(){
	utils.LoggingSettings(config.Config.LogFile)
	dao := implementations.PostDAOImpl{}
	timestamp := models.NewDateTimeStamp(time.RFC3339, time.RFC3339)
	post := models.NewPost(-1, "Title", "HELLO World",*timestamp)
	fmt.Println(dao.Create(*post))
}