package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/linebot"
)

func messageReply(event linebot.Event) (err error) {
	message := event.Message.(*linebot.TextMessage)

	if message.Text == "你在哪交換" {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewLocationMessage("捷克理工大學", "Zikova 1903/4, 166 36 Praha 6", 50.102974, 14.391177)).Do(); err != nil {
			log.Print(err)
			return err
		}
	} else if message.Text == "給我看履歷" {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("https://docs.google.com/document/d/1eKjTU4ebTTUJrISXR5YcQ47Jki-Br4apP_EQT_K_dbA/edit?usp=sharing")).Do(); err != nil {
			log.Print(err)
			return err
		}
	} else if message.Text == "你會哪些東西呢" {
		skills := []string{"Golang", ".NET MVC", "Gitlab CI/CD", "Mysql", "Redis"}
		replyMessage := "我會很多東西呢\n\n"

		for _, skill := range skills {
			replyMessage += "\n" + skill
		}
		replyMessage += "\n\n等等等..."

		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			log.Print(err)
			return err
		}
	} else if message.Text == "最近文章" {
		template := crawlBlog(4)
		packMessage := linebot.NewTemplateMessage("哎呀~ 這裡怎麼看不到呢", template)
		if _, err := bot.ReplyMessage(event.ReplyToken, packMessage).Do(); err != nil {
			log.Print(err)
			return err
		}
	} else if message.Text == "嗨" {

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

	} else {
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("哎呀！你這樣問我有點困擾呢！")).Do(); err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}

func crawlBlog(num int) *linebot.CarouselTemplate {
	template := linebot.NewCarouselTemplate()
	carouselCols := []*linebot.CarouselColumn{}
	// Request the HTML page.
	res, err := http.Get("https://a28283878.github.io/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".posts-wrapper article").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a").Text()
		postURL, _ := s.Find("a").Attr("href")
		pictureURL, _ := s.Find("a").Find("div").Attr("style")
		pictureURL = pictureURL[strings.Index(pictureURL, "(")+1 : strings.Index(pictureURL, ")")]

		btn := linebot.NewURITemplateAction("觀看", postURL)
		column := linebot.NewCarouselColumn(pictureURL, "文章", title, btn)

		carouselCols = append(carouselCols, column)
	})

	template = linebot.NewCarouselTemplate(carouselCols...)

	return template
}
