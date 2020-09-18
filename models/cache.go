package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"telegrambot/entities"
	"telegrambot/utilities"
	"time"
)

func GetTruthOrDareFromRedis(truth_or_dare string, telegramId int, redisClient *redis.Client) string {
	log.Println("telegramID: ", telegramId, "choice: ", truth_or_dare)
	key := fmt.Sprintf("user_%d", telegramId)
	data, err := utilities.ReadFromRedis(key, redisClient)
	if err != nil {
		return ""
	}
	var challenges []entities.TruthAndDare
	if err := json.Unmarshal([]byte(data), &challenges); err != nil {
		log.Println("Unable to convert json from redis to []entities.TruthAndDare.", err)
		return ""
	}

	length := len(challenges)

	if length< 1 {
		log.Println("No challenges in redis. Game over")
		return "Game over"
	}

	currentTimeNanoSeconds := time.Now().UnixNano()

	rand.Seed(currentTimeNanoSeconds)
	position := rand.Intn(length)
	challenge := challenges[position]

	//tries := 0
	//for challenge.Type != truth_or_dare && tries < 200  {
	//	position := rand.Intn(length)
	//	challenge = challenges[position]
	//	log.Println("== trial : ", tries," position: ", position, " type: ", challenge.Type, "choice: ", truth_or_dare)
	//	tries ++
	//}

	log.Println("The chosen challenge :: position ", position, "challenge: ", challenge.Challenge)
	// reference : https://yourbasic.org/golang/delete-element-slice/
	log.Println("The original array length is ", len(challenges))
	challenges[position] = challenges[length-1]
	challenges[length-1] = entities.TruthAndDare{}
	challenges = challenges[:length-1]
	log.Println("The new array length is ", len(challenges))

	_ = utilities.SaveToRedis(key, nil, redisClient)
	if err := utilities.SaveToRedis(key, challenges, redisClient); err != nil {
		log.Println("Unable to update redis because ", err)
	}

	return challenge.Challenge
}

