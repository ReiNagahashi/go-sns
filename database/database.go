package database

type Database interface {
    GetTableLength(tableName string)(int, error)
    PrepareAndExecute(query string, args ...interface{}) error
    PrepareAndFetchAll(query string, args ...interface{}) ([]map[string]interface{}, error)
}