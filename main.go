package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"telegrambot/entities"
	"telegrambot/resources"
	resources2 "telegrambot/resources/db"
	"telegrambot/resources/db/utils"
	"time"
)

type Message struct {
	CallbackQuery struct {
		From struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
		} `json:"from"`
		Data    string `json:"data"`
		Message struct {
			Chat struct {
				Id int `json:"id"`
			} `json:"chat"`
		} `json:"message"`
	} `json:"callback_query"`

	Message struct {
		Text string `json:"text"`
		From struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
		} `json:"from"`
		Chat struct {
			Id int `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func listen(w http.ResponseWriter, r *http.Request) {
	var data Message
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
	err = sendmessage(chatId, response, keyboard)
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
	initLogger()
	log.Println("Running")

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
		challenge := gettruthordare("truth", telegramId)
		if challenge == "" {
			return "An error occured, try again", resources.TruthOrDareKeyboard(gameId)
		}
		if challenge == "Game over" {
			return "Game Over.", ""
		}
		return challenge, resources.AcceptDeclineKeyboard(gameId)
	}
	if strings.Contains(message, "dare") {
		gameId := strings.Replace(message, "dare-", "", 1)
		challenge := gettruthordare("dare", telegramId)
		if challenge == "" {
			return "An error occured, try again", resources.TruthOrDareKeyboard(gameId)
		}
		if challenge == "Game over" {
			return "Game Over.", ""
		}
		return challenge, resources.AcceptDeclineKeyboard(gameId)
	}
	if message == "start" || message == "/start" {
		// TODO: When someone chooses Three, create a game session
		gameId, err := resources2.CreateGameSession()
		if err != nil {
			return "An error occured. Please try again later.", ""
		}
		return "Welcome to Truth or Dare game. How many players are you?", resources.PlayerCountKeyboard(gameId)
	}
	if strings.Contains(message, "players") {
		// e.g if the message is players3-6 it prints 3 players, game 6
		message = strings.Replace(message, "players", "", 1)
		arr := strings.Split(message, "-")
		numberOfPlayers := arr[0]
		gameId := arr[1]
		if err := resources2.SetGamePlayers(gameId, numberOfPlayers); err != nil {
			return "An error occured. Send 'start' to try again", ""
		}
		return fmt.Sprintf("Kindly tell your friends to text me 'Join %s'", gameId), ""
	}
	if strings.Contains(message, "join") {
		userId, err := resources2.Create_player(telegramId, firstName)
		if err != nil {
			return "an error occured", ""
		}
		gameId := strings.Replace(message, "join ", "", 1)
		gameId = strings.Trim(gameId, " ")

		// 1. Check if game session has enough players
		ok, err := resources2.SpaceAvailableInGameSession(gameId)
		if err != nil {
			return "an error occured", ""
		}
		if !ok {
			return "The game already has enough players.", ""
		}
		err = resources2.Scores(gameId, userId)
		if err != nil {
			return "an error occured", ""
		}
		key := fmt.Sprintf("user_%d", telegramId)
		truthsAndDares, err := utils.TruthsAndDaresFromDB()
		if err != nil {
			return "an error occured", ""
		}
		redisClient := utils.ConnectToRedis()
		err = utils.SaveToRedis(key, truthsAndDares, redisClient)
		if err != nil {
			return "an error occured", ""
		}

		return "congratulations Maxine for joining. Please choose truth or dare", resources.TruthOrDareKeyboard(gameId)
	}
	return "Please send 'start'", ""

}

func sendmessage(chatid int, message, keyboard string) (err error) {
	token := os.Getenv("TOKEN")
	msg := url.QueryEscape(message)
	link := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s&reply_markup=%s", token, chatid, msg, keyboard)
	log.Println("Sending message :: ", link)
	client := &http.Client{}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	log.Println("Message sent back to telegram")
	return nil
}

func initLogger() {
	logFolder := os.Getenv("LOG_FOLDER")
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s%s.log", logFolder+"app-", "%Y-%m-%d.%H%M"),
		rotatelogs.WithLinkName(logFolder+"link.log"),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(10000),
	)
	if err != nil {
		log.Println("Failed to initialize log file ", err.Error())
		log.SetOutput(os.Stderr)
		return
	}
	if os.Getenv("ENV") == "dev" {
		log.SetOutput(os.Stderr)
		return
	}
	log.SetOutput(writer)
	return
}

func gettruthordare(truth_or_dare string, telegramId int) string {
	log.Println("telegramID: ", telegramId, "choice: ", truth_or_dare)
	key := fmt.Sprintf("user_%d", telegramId)
	redisClient := utils.ConnectToRedis()
	data, err := utils.ReadFromRedis(key, redisClient)
	if err != nil {
		log.Println("An error occured. Please try again.", err)
		return ""
	}
	var challenges []entities.TruthAndDare
	if err := json.Unmarshal([]byte(data), &challenges); err != nil {
		log.Println("An error occured. Please try again.", err)
		return ""
	}

	length := len(challenges)

	if length< 1 {
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

	_ = utils.SaveToRedis(key, nil, redisClient)
	if err := utils.SaveToRedis(key, challenges, redisClient); err != nil {
		log.Println("Unable to update redis because ", err)
	}

	return challenge.Challenge
}