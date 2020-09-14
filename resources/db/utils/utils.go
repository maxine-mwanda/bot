package utils

import (
	"encoding/json"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
	"os"
	"telegrambot/entities"
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


/*func truth_or_Dares () {

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
}*/
func ConnectToRedis() (conn *redis.Client) {
	conn = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:3306",
			Password: "",
			DB:       0,
		},
	)
	return
}

func SaveToRedis(Key string, data interface{}, client *redis.Client) (err error) {
	expiryTime := time.Duration(time.Minute * 30)

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
