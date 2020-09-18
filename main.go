package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"telegrambot/entities"
	"telegrambot/models"
	"telegrambot/utilities"
)

var dbConn *sql.DB
var redisClient *redis.Client

func listen(w http.ResponseWriter, r *http.Request) {
	var data entities.Message
	var chatId int
	var msg string
	var firstName string
	var telegramId int

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Unable t decode payload because ", err)
		w.Write([]byte("error"))
		return
	}
	log.Println("Received payload")

	if data.CallbackQuery.Data == "" {
		// its a message
		chatId = data.Message.Chat.Id
		msg = data.Message.Text
		telegramId = data.Message.From.ID
		firstName = data.Message.From.FirstName
	} else {
		// its a callback (button, keyboard)
		chatId = data.CallbackQuery.Message.Chat.Id
		msg = data.CallbackQuery.Data
		telegramId = data.CallbackQuery.From.ID
		firstName = data.CallbackQuery.From.FirstName
	}

	log.Println("message", msg, "chat Id ", chatId)
	response, keyboard := getresponse(msg, firstName, telegramId)
	log.Println("response", response)
	log.Println("keyboard", keyboard)
	err = utilities.Sendmessage(chatId, response, keyboard)
	if err != nil {
		log.Println("error", err)
	}

	w.Write([]byte("ok"))
}

func main() {
	var port = ":3000"
	if err := godotenv.Load(); err != nil {
		log.Println("Unable to read .env file. exiting")
	}
	utilities.InitLogger()
	log.Println("Running")

	redisClient = utilities.ConnectToRedis()
	dbConn = utilities.Connecttodb()


	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "content-type", "content-length", "accept-encoding", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})

	router := mux.NewRouter()

	router.HandleFunc("/listen", listen).Methods("POST")

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(port, handlers.CORS(origins, headers, methods)(router)); err != nil {
		log.Printf("Unable to start API because %s", err.Error())
		os.Exit(3)
	}

}

func getresponse(message, firstName string, telegramId int) (string, string) {
	message = strings.ToLower(message)
	log.Println("The Message == ", message)

	if strings.Contains(message, "truth") {
		gameId := strings.Replace(message, "truth-", "", 1)
		challenge := models.GetTruthOrDareFromRedis("truth", telegramId, redisClient)
		if challenge == "" {
			return "An error occured, try again", utilities.TruthOrDareKeyboard(gameId)
		}
		if challenge == "Game over" {
			score, err := models.GetFinalScore(gameId, telegramId, dbConn)
			if err != nil {
				return "Game Over.", ""
			}
			return fmt.Sprintf("You have %d points.\nGame Over", score), ""
		}
		return challenge, utilities.AcceptDeclineKeyboard(gameId)
	}
	if strings.Contains(message, "dare") {
		gameId := strings.Replace(message, "dare-", "", 1)
		challenge := models.GetTruthOrDareFromRedis("dare", telegramId, redisClient)
		if challenge == "" {
			return "An error occured, try again", utilities.TruthOrDareKeyboard(gameId)
		}
		if challenge == "Game over" {
			score, err := models.GetFinalScore(gameId, telegramId, dbConn)
			if err != nil {
				return "Game Over.", ""
			}
			return fmt.Sprintf("You have %d points.\nGame Over", score), ""
		}
		return challenge, utilities.AcceptDeclineKeyboard(gameId)
	}
	if message == "start" {
		gameId, err := models.CreateGameSession(dbConn)
		if err != nil {
			return "An error occured. Please try again later.", ""
		}

		return "How many players are you?", utilities.PlayerCountKeyboard(gameId)
	}
	if message == "/start" {
		return "Welcome to truth or dare. \n1. Start the game by typing 'start'. \n2. Select the number of players. \n3. Join the game by sending the game session number. \n4. You have the option of accepting or declining a challenge. If you accept you earn ten points, if you decline you get 0.", ""
	}
	if strings.Contains(message, "accept") {
		gameId := strings.Replace(message, "accept-", "", 1)
		err := models.UpdatePlayerScore(telegramId, gameId, dbConn)
		if err != nil {
			return "An error occured, please try again", utilities.TruthOrDareKeyboard(gameId)
		}
		return "You have been awarded 10 points.", utilities.TruthOrDareKeyboard(gameId)
	}
	if strings.Contains(message, "decline") {
		gameId := strings.Replace(message, "decline-", "", 1)
		return "You have not been awarded any points.", utilities.TruthOrDareKeyboard(gameId)
	}


	if strings.Contains(message, "players") {
		message = strings.Replace(message, "players", "", 1)
		arr := strings.Split(message, "-")
		numberOfPlayers := arr[0]
		gameId := arr[1]
		if err := models.SetNumberOfGamePlayers(gameId, numberOfPlayers, dbConn); err != nil {
			return "An error occured. Send 'start' to try again", ""
		}
		return fmt.Sprintf("Please reply with 'Join %s'. Also, kindly tell your friends to text me 'Join %s'", gameId, gameId), ""
	}
	if strings.Contains(message, "join") {
		userId, err := models.CreatePlayer(telegramId, firstName, dbConn)
		if err != nil {
			return "an error occured", ""
		}
		gameId := strings.Replace(message, "join ", "", 1)
		gameId = strings.Trim(gameId, " ")

		// 1. Check if game session has enough players
		ok, err := models.SpaceAvailableInGameSession(gameId, dbConn)
		if err != nil {
			return "an error occured", ""
		}
		if !ok {
			return "The game already has enough players.", ""
		}
		err = models.Scores(gameId, userId, dbConn)
		if err != nil {
			return "an error occured", ""
		}
		key := fmt.Sprintf("user_%d", telegramId)
		truthsAndDares, err := models.TruthsAndDaresFromDB(dbConn)
		if err != nil {
			return "an error occured", ""
		}
		err = utilities.SaveToRedis(key, truthsAndDares, redisClient)
		if err != nil {
			return "an error occured", ""
		}

		return fmt.Sprintf("Congratulations %s for joining. Please choose truth or dare", firstName), utilities.TruthOrDareKeyboard(gameId)
	}
	return "Please send 'start'", ""

}
