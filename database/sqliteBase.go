package database

import (
	"database/sql"
	"fmt"
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

func (s *SqliteBase) PrepareAndFetchAll(query string, args ...interface{}) ([]map[string]interface{}, error){
	stmt, err := s.DbConnection.Prepare(query)
	if err != nil{
		return nil, err
	}
	defer stmt.Close()

	rows, execErr := stmt.Query(args...)
	if execErr != nil{
		return nil, execErr
	}
	defer rows.Close()
	
	columns, colErr := rows.Columns()
	if colErr != nil{
		return nil, colErr
	}

	results := []map[string]interface{}{}

	for rows.Next() {
		datas := make([]interface{}, len(columns))
		datasPtr := make([]interface{}, len(columns))

		for i := range datas{
			datasPtr[i] = &datas[i]
		}

		if scanErr := rows.Scan(datasPtr...); scanErr != nil{
			return nil, scanErr
		}

		rowMap := map[string]interface{}{}
		for i, col := range columns{
			rowMap[col] = datas[i]
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func (s *SqliteBase) PrepareAndExecute(query string, args ...interface{}) error {
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


func (s *SqliteBase) GetTableLength(tableName string)(int, error){
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var count int

	err := s.DbConnection.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}


func (s *SqliteBase) GetLastInsertedId()(int, error){
	var lastInsertId int
	err := s.DbConnection.QueryRow("SELECT last_insert_rowid()").Scan(&lastInsertId)
	if err != nil{
		return -1, err
	}
	
	return lastInsertId, nil
}