package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

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
