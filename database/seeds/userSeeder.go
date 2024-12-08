package seeds

import (
	"fmt"
	"go-sns/database"
	"log"
	"github.com/go-faker/faker/v4"
)

func UserSeed(SqliteBase *database.SqliteBase){
	stmt, err := SqliteBase.DbConnection.Prepare("INSERT INTO users (name, email, password) VALUES (?, ?, ?)")
	if err != nil{
		log.Fatal(err)
	}
	defer stmt.Close()

    numRecords := 10
    for i := 0; i < numRecords; i++{
        name := faker.Username()
        description := faker.Email()
		password := faker.Password()
        
        _, err := stmt.Exec(name, description, password)
        if err != nil{
            log.Fatal(err)
        }
    }

    fmt.Println("Seeding Completed")
}