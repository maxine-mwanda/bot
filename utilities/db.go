package utilities

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func  Connecttodb() (connection *sql.DB) {
	dburi := os.Getenv("DBURI")
	connection, err := sql.Open("mysql", dburi)
	if err != nil {
		log.Println("unable to connect to db", err)
		os.Exit(3)
	}
	log.Println("Connected to db successfully")
	connection.SetMaxOpenConns(100)
	connection.SetMaxIdleConns(10)
	connection.SetConnMaxIdleTime(time.Second * 10)

	return
}


