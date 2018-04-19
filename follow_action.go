package main

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)

func followAction(event linebot.Event) (err error) {

	locationBtn := linebot.NewMessageTemplateAction("你在哪交換", "你在哪交換")
	resumeBtn := linebot.NewMessageTemplateAction("給我看履歷", "給我看履歷")
	skillBtn := linebot.NewMessageTemplateAction("你會哪些東西呢", "你會哪些東西呢")
	travelBtn := linebot.NewMessageTemplateAction("最近文章", "最近文章")

	template := linebot.NewButtonsTemplate("https://farm1.staticflickr.com/799/41548719091_313673967f_b.jpg", "這裡有些範例問題呢",
		"選個看看吧", locationBtn, resumeBtn, skillBtn, travelBtn)
	packMessage := linebot.NewTemplateMessage("哎呀~ 這裡怎麼看不到呢", template)
	if _, err := bot.ReplyMessage(event.ReplyToken, packMessage).Do(); err != nil {
		log.Print(err)
		return err
	}

	return nil
}
