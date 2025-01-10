package utils

import (
	"database/sql"
	"go-sns/config"
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

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil{
		log.Fatalf("failed to enable foreign keys: %v", err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil{
		log.Fatalf("failed to set dialect: %v", err)
	}

	migrationsDir := "./database/migrations"
	if err := goose.Up(db, migrationsDir); err != nil{
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")

}

