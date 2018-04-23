package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var channelSecret string
var channelAccessToken string
var questionTemplate *linebot.ButtonsTemplate

func main() {
	var err error

	initValues()
	bot, err = linebot.New(channelSecret, channelAccessToken)
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func initValues() {
	channelSecret = os.Getenv("ChannelSecret")
	channelAccessToken = os.Getenv("ChannelAccessToken")

	locationBtn := linebot.NewMessageTemplateAction("你在哪交換", "你在哪交換")
	skillBtn := linebot.NewMessageTemplateAction("你會哪些東西呢", "你會哪些東西呢")
	hobbyBtn := linebot.NewMessageTemplateAction("你有哪些興趣呢", "你有哪些興趣呢")
	travelBtn := linebot.NewMessageTemplateAction("最近文章", "最近文章")

	questionTemplate = linebot.NewButtonsTemplate("https://farm1.staticflickr.com/799/41548719091_313673967f_b.jpg", "這裡有些範例問題呢",
		"選個看看吧", locationBtn, skillBtn, hobbyBtn, travelBtn)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	//parse post
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
			log.Print(err)
		} else {
			w.WriteHeader(500)
			log.Print(err)
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
					log.Print(err)
					w.WriteHeader(500)
					return
				}
			default:
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("哎呀！沒有辦法回答這東西呢！")).Do(); err != nil {
					log.Print(err)
					return
				}
				log.Printf("Unhandle message type : %v", &event)
			}
		} else if event.Type == linebot.EventTypeFollow {
			err = followAction(event)
			if err != nil {
				w.WriteHeader(500)
				log.Print(err)
				return
			}
		} else {
			log.Printf("Unhandle event type : %v", &event)
			if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("哎呀！沒有辦法回答這東西呢！")).Do(); err != nil {
				log.Print(err)
				return
			}
		}
	}
}
