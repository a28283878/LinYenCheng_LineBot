package main

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

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
	} else if message.Text == "嗨" {
		locationBtn := linebot.NewMessageTemplateAction("你在哪交換", "你在哪交換")
		resumeBtn := linebot.NewMessageTemplateAction("給我看履歷", "給我看履歷")

		template := linebot.NewConfirmTemplate("這裡有範例問題可以按歐~", locationBtn, resumeBtn)
		message := linebot.NewTemplateMessage("哎呀~ 這裡怎麼看不到呢", template)
		if _, err := bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
			log.Print(err)
		}
	} else {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("哎呀！你這樣問我有點困擾呢！")).Do(); err != nil {
			log.Print(err)
		}
	}

}
