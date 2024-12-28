package utils

import (
	"database/sql"
	"go-sns/config"
	"go-sns/database"
	"go-sns/database/seeds"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLiteドライバ
	"github.com/pressly/goose/v3"
)

func RunMigrations(){
    db, err := sql.Open(config.Config.SQLDriver, config.Config.DbPath)
	if err != nil{
		log.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()

	if err := goose.SetDialect("sqlite3"); err != nil{
		log.Fatalf("failed to set dialect: %v", err)
	}

	migrationsDir := "./database/migrations"
	if err := goose.Up(db, migrationsDir); err != nil{
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")

}

func RunSeedings(){
	db := database.NewSqliteBase()
	defer db.DbConnection.Close()
	seeds.PostSeed(db)
	seeds.UserSeed(db)
}