package utils

import (
	"encoding/json"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

func truth_or_Dares () {

	for truths := 0; truths < 5; truths++ {
		if !(truths >= 5) {
			fmt.Println("answer question %d", truths)
			continue
		}
	}
	for dares := 0; dares < 5; dares++ {
		if !(dares >= 5) {
			fmt.Println("answer question %d", dares)
			continue
		}
	}
	_= connectToRedis()
}
func connectToRedis() (conn *redis.Client) {
	conn = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	)
	return
}

func SaveToRedis(Key string, data interface{}, client *redis.Client) (err error) {
	expiryTime := time.Duration(time.Minute * 2)

	jsonData, _ := json.Marshal(data)

	if err = client.Set(Key, jsonData, expiryTime).Err(); err != nil {
		log.Println("Unable to save data to redis because ", err)
		return
	}
	return
}

func ReadFromRedis(Key string, client *redis.Client) (data string, err error) {
	data, err = client.Get(Key).Result()
	if err != nil {
		fmt.Println("Unable to read ", Key, " because ", err)
		return
	}
	return
}
