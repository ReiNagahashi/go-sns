package seeds

import (
	"fmt"
	"go-sns/database"
	"log"
	"github.com/go-faker/faker/v4"
)

func PostSeed(sqlite *database.SqliteBase){
	stmt, err := sqlite.DbConnection.Prepare("INSERT INTO posts (title, description, submitted_by) VALUES (?, ?, ?)")
	if err != nil{
		log.Fatal(err)
	}
	defer stmt.Close()

    numRecords := 10
    table_len, err := sqlite.GetTableLength("posts")
    if err != nil{
        log.Fatalln(err)
    }
    for i := 0; i < numRecords; i++{
        title := faker.Paragraph()
        description := faker.Paragraph()
        submitted_by, err := faker.RandomInt(table_len)
        if err != nil{
            log.Fatalln(err)
        }
        
        _, err = stmt.Exec(title, description, submitted_by[0])
        if err != nil{
            log.Fatal(err)
        }
    }

    fmt.Println("Post Seeding Completed")
}