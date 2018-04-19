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
				err = messageReply(event)
				if err != nil {
					w.WriteHeader(500)
					return
				}
			}
		} else if event.Type == linebot.EventTypeFollow {
			err = followAction(event)
			if err != nil {
				w.WriteHeader(500)
				return
			}
		} else {
			if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("哎呀！沒有辦法回答這東西呢！")).Do(); err != nil {
				log.Print(err)
				return
			}
		}
	}
}
