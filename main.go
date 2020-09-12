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
	"telegrambot/resources"
	resources2 "telegrambot/resources/db"
	"time"
)

type Message struct {
	CallbackQuery struct {
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
		//telegramId = data.Message.From.ID
		//firstName = data.Message.From.FirstName
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
		fmt.Println("Unable to read .env file. exiting")
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

	if message == "truth" {
		return getTruth(), resources.AcceptDeclineKeyboard()
	}
	if message == "dare" {
		return getDare(), resources.AcceptDeclineKeyboard()
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
		// TODO: create a function that prints the game id and number of players.
		// e.g if the message is players3-6 it prints 3 players, game 6
		if strings.Contains(message, "players 3") {
			resources.PlayerCountKeyboard(6)
		}
		return "Kindly tell your friends to text me 'Join 567'", ""
	}
	if message == "join 567" {
		// TODO: When a player sends join 567, add the record to player_scores table
		resources2.Create_player(telegramId, firstName)

		return "congratulations Maxine for joining. Please choose truth or dare", resources.TruthOrDareKeyboard()
	}
	return "Please choose truth or dare", resources.TruthOrDareKeyboard()

}

func getTruth() string {

	currentTimeNanoSeconds := time.Now().UnixNano()
	rand.Seed(currentTimeNanoSeconds)

	var truths = [5]string{
		"When was the last time you lied?",
		"When was the last time you cried?",
		"What's your biggest fear?",
		"What's your biggest fantasy?",
		"Do you have any fetishes?",
	}
	position := rand.Intn(5)

	return truths[position]

	}

func getDare() string {

	currentTimeNanoSeconds := time.Now().UnixNano()
	rand.Seed(currentTimeNanoSeconds)

	var dares = [5]string{
		"Kiss the person to your left",
		"Attempt to do a magic trick",
		"Do four cartwheels in row",
		"Let someone shave part of your body",
		"Eat five tablespoons of a condiment",
	}
	position := rand.Intn(5)
	return dares[position]
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
		fmt.Println("Failed to initialize log file ", err.Error())
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
