package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = "I@wanna9873"
	dbname   = "godemo"
)

// Database instance
var DB *sql.DB

func GetConnection(connString string) error {
	var err error
	DB, err = sql.Open("mysql", connString)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	fmt.Println("Connection Opened to Database")

	return nil
}

func GetConnectionString() string {
	psqlInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbname)

	return psqlInfo
}
