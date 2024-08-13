package database

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

func ConnectMySQL(host, user, password, dbname, port string) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
                      user, password, host, port, dbname)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}
