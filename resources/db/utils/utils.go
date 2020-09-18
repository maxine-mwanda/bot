package utils

import (
	"encoding/json"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"os"
	"telegrambot/entities"
	"time"

)

func  Connecttodb() (connection *sql.DB, err error) {
	dburi := os.Getenv("DBURI")
	connection, err = sql.Open("mysql", dburi)
	if err != nil {
		log.Println("unable to connect to db", err)
		os.Exit(3)
	}
	return
}


func TruthsAndDaresFromDB() (truthsAndDares []entities.TruthAndDare , err error) {
	db, err := Connecttodb()
	if err != nil {
		log.Println("unable to connect todb")
		return
	}
	query := "select * from truths_dares;"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("unable to fetch truths or dares")
		return
	}

	var t entities.TruthAndDare
	for rows.Next(){
		err = rows.Scan(&t.ID, &t.Challenge, &t.Type)
		if err != nil {
			log.Println("unable to scan from db")
			continue
		}
		truthsAndDares = append(truthsAndDares, t)
	}
	return
}


func ConnectToRedis() (conn *redis.Client) {
	conn = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	)

	if err := conn.Ping().Err(); err != nil {
		log.Println("Unable to connect to redis. Exiting...", err)
		os.Exit(3)
	}
	return
}

func SaveToRedis(Key string, data interface{}, client *redis.Client) (err error) {
	expiryTime := time.Duration(time.Minute * 30)

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Unable to convert data to json because", err.Error())
		return
	}

	if err = client.Set(Key, jsonData, expiryTime).Err(); err != nil {
		log.Println("Unable to save data to redis because ", err.Error())
		return
	}
	log.Println("Data saved to redis key ", Key)
	return
}

func ReadFromRedis(Key string, client *redis.Client) (data string, err error) {
	data, err = client.Get(Key).Result()
	if err != nil {
		log.Println("Unable to read ", Key, " because ", err)
		return
	}
	return
}

