package seeds

import (
	"go-sns/database"
	"log"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
)

func PostSeed(sqlite *database.SqliteBase, wg *sync.WaitGroup){
    defer wg.Done()
	stmt, err := sqlite.DbConnection.Prepare("INSERT INTO posts (title, description, submitted_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
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
        now := time.Now()
        _, err = stmt.Exec(title, description, submitted_by[0], now, now)
        if err != nil{
            log.Fatal(err)
        }
    }

}