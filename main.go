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

	timestamp := models.NewDateTimeStamp(time.Now(), time.Now())
	post := models.NewPost(-1, "Title", "HELLO World",*timestamp)
	fmt.Println(dao.Create(*post))
	posts := dao.GetAll(0, 20)
	for _, post := range posts{
		fmt.Println(post.GetTimeStamp().GetCreatedAt())
	}
	// db := database.NewSqliteBase()
	// seeds.PostSeed(db)
	// routing.StartWebServer()
}