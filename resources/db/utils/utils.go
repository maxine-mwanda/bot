package utils

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

func  Connecttodb() (connection *sql.DB, err error) {
	dburi := "root:daddy@tcp(localhost:3306)/Truth_or_dare"
	connection, err = sql.Open("mysql", dburi)
	if err != nil {
		fmt.Println("unable to connect to db", err)
		return
	}
	return
}
