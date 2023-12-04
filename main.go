package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pixie/api"
	"pixie/api/sentence"
	"pixie/api/user"
	"pixie/api/vocabulary"
	"pixie/core/pixie"
	"pixie/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sashabaranov/go-openai"
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

func CorsConfig() cors.Config {
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"}
	corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	return corsConf
}

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

	// {
	// 	upgrade()
	// }

	bot, err := linebot.New(lineBotChannelSecret, lineBotChannelToken)
	if err != nil {
		panic(err)
	}
	LineBot = bot

	pixie.Build(openAIToken)

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./out", false)))
	router.GET("/", func(ctx *gin.Context) {
		ctx.File("./out/index.html")
	})

	router.Use(cors.New(CorsConfig()))

	router.POST("/callback", callbackHandler)
	router.POST("/api/user/member/login", user.Member.Login)

	logined := router.Group("", api.CheckAuth)
	{
		logined.POST("/api/vocabulary/pool/ask", vocabulary.Pool.Ask)
		logined.POST("/api/vocabulary/pool/back", vocabulary.Pool.Back)
		logined.POST("/api/vocabulary/pool/forward", vocabulary.Pool.Forward)
		logined.POST("/api/vocabulary/pool/jump_to", vocabulary.Pool.JumpTo)
		logined.POST("/api/vocabulary/bookmark/toggle", vocabulary.Bookmark.Toggle)
		logined.POST("/api/vocabulary/bookmark/browse", vocabulary.Bookmark.Browse)
		logined.POST("/api/vocabulary/unknown/ask", vocabulary.Unknown.Ask)

		logined.POST("/api/sentence/interpret", sentence.Interpret)
		logined.POST("/api/sentence/explain", sentence.Explain)

	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	fmt.Println("service stop ....")
	// http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("1"))
	// })
	// http.HandleFunc("/callback", callbackHandler)

	// http.ListenAndServe(":8080", nil)
}

func callbackHandler(ctx *gin.Context) {
	// ctx := context.Background()

	events, err := LineBot.ParseRequest(ctx.Request)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			ctx.JSON(400, nil)
		} else {
			ctx.JSON(500, nil)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				reply, err := pixie.Chat([]openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: message.Text,
					},
				})
				if err != nil {
					if _, err = LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do(); err != nil {
						log.Print(err)
					}
					return
				}

				if _, err = LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
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
