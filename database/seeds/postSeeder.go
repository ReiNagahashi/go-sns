package seeds

import (
	"fmt"
	"go-sns/database"
	"log"
	"github.com/go-faker/faker/v4"
)

func PostSeed(SqliteBase *database.SqliteBase){
	stmt, err := SqliteBase.DbConnection.Prepare("INSERT INTO posts (title, description) VALUES (?, ?)")
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