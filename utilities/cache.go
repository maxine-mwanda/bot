package utilities

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
)

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
	log.Println("Data ", Key, "read successfully from redis")
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

