package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	CallbackQuery struct {
		Data string `json:"data"`
		Message struct{
			Chat struct{
				Id int `json:"id"`
			} `json:"chat"`
		} `json:"message"`
	} `json:"callback_query"`

	Message struct {
		Text string `json:"text"`
		Chat struct {
			Id int `json:"id"`
			FirstName string `json:"first_name"`
		} `json:"chat"`
	} `json:"message"`
}

func listen(w http.ResponseWriter, r *http.Request) {
	var data Message
	var chatId int
	var msg string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.Write([]byte("error"))
		return

	}

	if data.CallbackQuery.Data == "" {
		// its a message
		chatId = data.Message.Chat.Id
		msg = data.Message.Text
	} else {
		// its a callback
		chatId = data.CallbackQuery.Message.Chat.Id
		msg = data.CallbackQuery.Data
	}

	fmt.Println("message", msg)
	response := getresponse(msg)
	fmt.Println("response", response)
	keyboard:= CreateKeyboard()
	err= sendmessage(chatId, response, keyboard)
	if err != nil {
		fmt.Println("error", err)
	}

	w.Write([]byte("ok"))
}

func main() {
	var port = ":3000"
	_ = godotenv.Load()


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


func getresponse(message string) string{
	message = strings.ToLower(message)

	switch message {
	case "truth" :
		return "Tell me your secret"
	case "dare" :
		return "Scream at someone"
	default:
		return "Please choose truth or dare"

	}
}

func sendmessage (chatid int, message, keyboard string) (err error) {
	token := os.Getenv("TOKEN")
	url := fmt.Sprintf( "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s&reply_markup=%s", token, chatid, message, keyboard)
	client := &http.Client{}

	req,err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode!=200{
		return errors.New(res.Status)
	}
	return nil
}

func CreateKeyboard () string {
	keyboard:= map[string]interface{}{
		"inline_keyboard" : [][]map[string]string{
			{
			{
				"text": "Truth",
				"callback_data": "truth",
			},
			},

			{
				{
					"text": "Dare",
					"callback_data": "dare",
				},
			},
		},
	}

	jsonkeyboard, _ := json.Marshal(keyboard)
	return string(jsonkeyboard)
}