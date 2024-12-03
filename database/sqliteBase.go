package database

import (
	"database/sql"
	"fmt"
	"go-sns/config"
	"log"

	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

func init() {
	var err error
	DbConnection, err = sql.Open(config.Config.SQLDriver, config.Config.DbPath)
	if err != nil{
		log.Fatalln(err)
	}
	defer DbConnection.Close()

	stmt, err := DbConnection.Prepare("INSERT INTO posts (title, description) VALUES (?, ?)")
	if err != nil{
		log.Fatal(err)
	}
	defer stmt.Close()

	numRecords := 10
	for i := 0; i < numRecords; i++{
		title := faker.Paragraph()
		description := faker.Paragraph()
		
		_, err := stmt.Exec(title, description)
		if err != nil{
			log.Fatal(err)
		}
	}

	fmt.Println("Seeding Completed")

}





