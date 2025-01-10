package main

import (
	"fmt"
	"go-sns/config"
	"go-sns/database"
	"go-sns/database/seeds"
	"go-sns/routing"
	"go-sns/utils"
	"log"
	"sync"
	"time"
)

func main(){
	utils.LoggingSettings(config.Config.LogFile)
	utils.RunMigrations()

	db := database.NewSqliteBase()
	defer db.DbConnection.Close()
	var wg sync.WaitGroup

	go func(){
		for{
			wg.Add(1)
			time.Sleep(60 * time.Second)
			go seeds.PostSeed(db, &wg)
			wg.Wait()
			fmt.Println("Post Seeding Completed")
		}
	}()

	log.Println(routing.StartWebServer())
}
