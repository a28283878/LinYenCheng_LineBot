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
				imageURL := "/assests/img/YDXJ0571.jpg"
				template := linebot.NewCarouselTemplate(
					linebot.NewCarouselColumn(
						imageURL, "hoge", "fuga",
						linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
						linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
					),
					linebot.NewCarouselColumn(
						imageURL, "hoge", "fuga",
						linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは"),
						linebot.NewMessageTemplateAction("Say message", "Rice=米"),
					),
				)

				if _, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTemplateMessage("Carousel alt text", template),
				).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
