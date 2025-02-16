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

	// "net/http"
	_ "net/http/pprof"
)

// TODO: postsからcreated_atをもとに高順にリミットの数だけ取得するところがうまく動作しない

func main() {
	// go func() {
    //     log.Println(http.ListenAndServe("localhost:6060", nil))
    // }()

	utils.LoggingSettings(config.Config.LogFile)
	utils.RunMigrations()
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()
	var wg sync.WaitGroup
	// seeds.UserSeed(db)
	go func(){
		for{
			wg.Add(2)
			time.Sleep(10000 * time.Second)
			go seeds.PostSeed(db, &wg)
			go seeds.PostFavoriteSeed(db, &wg)
			wg.Wait()
			fmt.Println("Seeding Completed")
		}
	}()

	log.Println(routing.StartWebServer())
}
