package main

import (
	"commuting-time-line-bot/handler"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.Print("[Start] start main")
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	log.Printf("port : %s", port)
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// 実際にRequestを受け取った時に処理を行うHandle関数を定義し、handlerに登録
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Print("[Start] start handle events")
		defer log.Print("[End] end handle events")

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

		for _, event := range events {
			if event.Type != linebot.EventTypeMessage {
				return
			}

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// 飛んできたメッセージが正しい形式化を確認
				if err := strings.Contains(message.Text, "から"); err != true {
					replyMessage := "経路のお願いの仕方が違うぞ！\n （駅名）から（駅名）でたのむ！！"
					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Printf("contains : %s", err)
					}
				}
				// 出発駅と到着駅を分割
				slice := strings.Split(message.Text, "から")
				stations := map[string]string{}
				key := [...]string{"startStation", "targetStation"}
				log.Printf("message from line : %v", slice)
				for i, station := range slice {
					stations[key[i]] = station
				}

				// 経路情報のページへアクセスするurl
				routeUrl, err := handler.GetUrl(stations[key[0]], stations[key[1]])
				if err != nil {
					log.Printf("failed create url for route page : %s", err.Error())
				}

				// replyメッセージの作成
				rowReplyMessage := "こんな経路が見つかったぞ！\n" + routeUrl
				replyMessage := linebot.NewTextMessage(rowReplyMessage)
				log.Print(replyMessage)

				// 返信の実行
				if _, err = bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	})

	// /callback にエンドポイントの定義
	// HTTPサーバの起動
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}


