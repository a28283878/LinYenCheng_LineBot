package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var channelSecret string
var channelAccessToken string

func main() {
	var err error

	channelSecret = os.Getenv("ChannelSecret")
	channelAccessToken = os.Getenv("ChannelAccessToken")

	bot, err = linebot.New(channelSecret, channelAccessToken)
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	//valid signature
	decoded, err := base64.StdEncoding.DecodeString(r.Header.Get("X-Line-Signature"))
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	if !hmac.Equal(decoded, hash.Sum(nil)) {
		log.Printf("not post from Line server : %v", r)
		w.WriteHeader(400)
	}

	//parse post
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	//parse events
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch event.Message.(type) {
			case *linebot.TextMessage:
				messageReply(event)
			}
		}
	}
}

func messageReply(event linebot.Event) {
	message := event.Message.(*linebot.TextMessage)

	if message.Text == "你在哪交換" {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewLocationMessage("捷克理工大學", "Zikova 1903/4, 166 36 Praha 6", 50.102974, 14.391177)).Do(); err != nil {
			log.Print(err)
		}
	} else if message.Text == "給我看履歷" {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("https://docs.google.com/document/d/1eKjTU4ebTTUJrISXR5YcQ47Jki-Br4apP_EQT_K_dbA/edit?usp=sharing")).Do(); err != nil {
			log.Print(err)
		}
	} else {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("哎呀！你這樣問我有點困擾呢！")).Do(); err != nil {
			log.Print(err)
		}
	}

}
