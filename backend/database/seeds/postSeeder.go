package seeds

import (
	"fmt"
	"go-sns/database"

	"go-sns/database/dataAccess/implementations/userImpl"
	"go-sns/database/dataAccess/interfaces"
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

    var dao interfaces.UserDAO = userImpl.NewUserDAOImpl(sqlite)
    fmt.Printf("%v", dao)

    users, err := dao.GetAll()
    if err != nil{
        log.Fatalln(err)
    }

    if len(users) == 0{
        return
    }

    user_ids := make([]int, len(users))
    for i := range user_ids{
        user_ids[i] = users[i].GetId()
    }

    numRecords := 5

    for i := 0; i < numRecords; i++{
        title := faker.Paragraph()
        description := faker.Paragraph()
        randomInt, err := faker.RandomInt(0, len(users)-1)
        if err != nil{
            log.Fatalln(err)
        }
        submitted_by := user_ids[randomInt[0]]
        if err != nil{
            log.Fatalln(err)
        }
        now := time.Now()
        _, err = stmt.Exec(title, description, submitted_by, now, now)
        if err != nil{
            log.Fatal(err)
        }
    }

}