package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/linebot"
)

var postNumber = 4

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
	} else if message.Text == "你有哪些興趣呢" {
		skills := []string{"游泳", "打羽毛球", "看美劇", "踏青", "在床上耍廢"}
		replyMessage := "我有很多興趣呢\n\n"

		for _, skill := range skills {
			replyMessage += "\n" + skill
		}
		replyMessage += "\n\n等等等..."

		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			log.Print(err)
			return err
		}
	} else if message.Text == "最近文章" {
		template := crawlBlog(postNumber)
		packMessage := linebot.NewTemplateMessage("哎呀~ 要用手機看歐", template)
		if _, err := bot.ReplyMessage(event.ReplyToken, packMessage).Do(); err != nil {
			log.Print(err)
			return err
		}
	} else if message.Text == "嗨" {
		packMessage := linebot.NewTemplateMessage("哎呀~ 要用手機看歐", questionTemplate)
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
		if i >= num {
			return
		}
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
