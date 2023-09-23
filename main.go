package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"pixie/agent"
	"pixie/db"
	"pixie/pixie"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// func main() {
// 	ctx := context.Background()

// 	if err := db.Build(ctx, "cluster0.bi4tgm4.mongodb.net", "root", "DeiQ2mMvZmJwS9WZ"); err != nil {
// 		panic(err)
// 	}

// 	if err := db.Pixie().Collection(model.CAgent).Insert(ctx, model.Agent{
// 		UserId: "test",
// 	}); err != nil {
// 		panic(err)
// 	}
// 	// pixie.Build("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")
// 	// userID := "Leo"
// 	// log.Println("ready...")
// 	// for {
// 	// 	var commandData string
// 	// 	fmt.Scanln(&commandData)
// 	// 	message := agent.ExecuteCommand(userID, commandData)
// 	// 	log.Println(message.Marshal())
// 	// }
// }

var LineBot *linebot.Client = nil

func main() {
	dbDomain := os.Getenv("DB_DOMAIN")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	lineBotChannelSecret := os.Getenv("LINEBOT_CHANNEL_SECRET")
	lineBotChannelToken := os.Getenv("LINEBOT_CHANNEL_TOKEN")
	openAIToken := os.Getenv("OPENAI_TOKEN")

	fmt.Println("lineBotChannelSecret:", lineBotChannelSecret)
	fmt.Println("lineBotChannelToken:", lineBotChannelToken)
	fmt.Println("openAIToken:", openAIToken)

	if err := db.Build(context.Background(), dbDomain, dbUser, dbPassword); err != nil {
		panic(err)
	}

	bot, err := linebot.New(lineBotChannelSecret, lineBotChannelToken)
	if err != nil {
		panic(err)
	}
	LineBot = bot

	pixie.Build(openAIToken)
	// OpenAI = openai.NewClient(openAIToken)

	http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("1"))
	})
	http.HandleFunc("/callback", callbackHandler)

	http.ListenAndServe(":8080", nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	events, err := LineBot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
	quota, err := LineBot.GetMessageQuota().Do()
	if err != nil {
		log.Println("Quota err:", err)
	}
	if quota.Value <= 0 {
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			fmt.Println("event message type:", event.Message.Type())
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				replyMessage := agent.ExecuteCommand(ctx, event.Source.UserID, message.Text)

				if _, err = LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage.Marshal())).Do(); err != nil {
					log.Print(err)
				}
			default:
				fmt.Printf("receive unsupport messge: %+v\n", event.Message)
				if _, err = LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("not support this message type now")).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}

	return
}
