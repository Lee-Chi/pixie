package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"pixie/agent"
	"pixie/db"
	"pixie/pixie"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Config struct {
	DBDomain   string `json:"mongodb_domain"`
	DBUser     string `json:"mongodb_user"`
	DBPassword string `json:"mongodb_password"`

	LineBotChannelSecret string `json:"linebot_channel_secret"`
	LineBotChannelToken  string `json:"linebot_channel_token"`

	OpenAIToken string `json:"openai_token"`
}

var LineBot *linebot.Client = nil

func main() {
	dbDomain := os.Getenv("MONGODB_DOMAIN")
	dbUser := os.Getenv("MONGODB_USER")
	dbPassword := os.Getenv("MONGODB_PASSWORD")
	lineBotChannelSecret := os.Getenv("LINEBOT_CHANNEL_SECRET")
	lineBotChannelToken := os.Getenv("LINEBOT_CHANNEL_TOKEN")
	openAIToken := os.Getenv("OPENAI_TOKEN")

	conf := flag.String("config", "", "")
	flag.Parse()

	if *conf != "" {
		data, err := os.ReadFile(*conf)
		if err != nil {
			panic(err)
		}

		config := Config{}
		if err := json.Unmarshal(data, &config); err != nil {
			panic(err)
		}

		dbDomain = config.DBDomain
		dbUser = config.DBUser
		dbPassword = config.DBPassword

		lineBotChannelSecret = config.LineBotChannelSecret
		lineBotChannelToken = config.LineBotChannelToken

		openAIToken = config.OpenAIToken
	}

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
