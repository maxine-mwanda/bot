package utilities

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func Sendmessage(chatid int, message, keyboard string) (err error) {
	token := os.Getenv("TOKEN")
	msg := url.QueryEscape(message)
	link := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s&reply_markup=%s", token, chatid, msg, keyboard)
	log.Println("Sending message :: ", link)
	client := &http.Client{}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Println("Unable to create http request because ", err)
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Unable to make http request because ", err)
		return err
	}

	if res.StatusCode != 200 {
		log.Println(res.Status)
		return errors.New(res.Status)
	}
	log.Println("Message sent back to telegram")
	return nil
}
