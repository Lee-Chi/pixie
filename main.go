package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"pixie/core/pixie"

	"github.com/Lee-Chi/go-sdk/logger"
	"github.com/Lee-Chi/go-sdk/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func CorsConfig() cors.Config {
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"}
	corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	return corsConf
}

func main() {
	// dbDomain := os.Getenv("MONGODB_DOMAIN")
	// dbUser := os.Getenv("MONGODB_USER")
	// dbPassword := os.Getenv("MONGODB_PASSWORD")
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

		// dbDomain = config.DBDomain
		// dbUser = config.DBUser
		// dbPassword = config.DBPassword

		lineBotChannelSecret = config.LineBotChannelSecret
		lineBotChannelToken = config.LineBotChannelToken

		openAIToken = config.OpenAIToken
	}

	logger.Init()

	// if err := db.Build(context.Background(), dbDomain, dbUser, dbPassword); err != nil {
	// 	panic(err)
	// }

	// {
	// 	upgrade()
	// }

	bot, err := linebot.New(lineBotChannelSecret, lineBotChannelToken)
	if err != nil {
		panic(err)
	}
	LineBot = bot

	pixie.Build(openAIToken)

	logger.Info("server start")

	engine := gin.Default()
	router := engine.Group("", func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		id := uuid.New().String()
		service.Accept(id)
		defer func() {
			duration := service.Done(id)
			fmt.Printf("[REQUEST] | %vs | %s\n", float64(duration)/float64(time.Second), path)
		}()
		ctx.Next()
	})
	// router.Use(static.Serve("/", static.LocalFile("./out", false)))
	// router.GET("/", func(ctx *gin.Context) {
	// 	ctx.File("./out/index.html")
	// })

	// router.Use(cors.New(CorsConfig()))

	router.GET("/sleep", func(ctx *gin.Context) {
		time.Sleep(2 * time.Second)
		ctx.JSON(http.StatusOK, nil)
	})

	router.POST("/callback", callbackHandler)

	go func() {
		if err := engine.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	service.Wait(context.Background())

	// http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("1"))
	// })
	// http.HandleFunc("/callback", callbackHandler)

	// http.ListenAndServe(":8080", nil)
	logger.Info("server stop")
}

func callbackHandler(ctx *gin.Context) {
	events, err := LineBot.ParseRequest(ctx.Request)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			ctx.JSON(http.StatusBadRequest, nil)
		} else {
			ctx.JSON(http.StatusInternalServerError, nil)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				response, err := pixie.Bot.Chat(message.Text)
				if err != nil {
					logger.Error("pixie bot chat, %v", err)
					response = "I'm sorry, I can't do that"
				}

				if _, err = LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
					logger.Error("linebot reply message, %v", err)
					return
				}
			default:
				logger.Error("receive unsupport messge, %+v", event.Message)
				if _, err = LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("not support this message type now")).Do(); err != nil {
					logger.Error("linebot reply message, %v", err)
					return
				}
			}
		}
	}

	return
}
