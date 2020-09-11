package utils

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"os"
)

func  Connecttodb() (connection *sql.DB, err error) {
	dburi := os.Getenv("DBURI")
	connection, err = sql.Open("mysql", dburi)
	if err != nil {
		fmt.Println("unable to connect to db", err)
		return
	}
	return
}
