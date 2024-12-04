package database

import (
	"database/sql"
	"go-sns/config"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteBase struct{
	DbConnection *sql.DB
}

func NewSqliteBase() *SqliteBase{
	var err error
	SqliteBase := &SqliteBase{}

	SqliteBase.DbConnection, err = sql.Open(config.Config.SQLDriver, config.Config.DbPath)
	if err != nil{
		log.Fatalln(err)
	}

	return SqliteBase
}

func (s *SqliteBase) Execute(query string, args ...interface{}) error {
	stmt, err := s.DbConnection.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, execErr := stmt.Exec(args...)
	if execErr != nil {
		return execErr
	}

	return nil
}

