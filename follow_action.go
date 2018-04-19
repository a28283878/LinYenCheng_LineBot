package main

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func followAction(event linebot.Event) (err error) {
	packMessage := linebot.NewTemplateMessage("哎呀~ 這裡怎麼看不到呢", questionTemplate)
	if _, err := bot.ReplyMessage(event.ReplyToken, packMessage).Do(); err != nil {
		log.Print(err)
		return err
	}

	return nil
}
