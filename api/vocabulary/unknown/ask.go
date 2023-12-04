package unknown

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pixie/api"
	"pixie/core/pixie"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type Answer struct {
	Word string `json:"word"`
	Vocs []Voc  `json:"vocs"`
}

type Voc struct {
	Word         string   `json:"word"`
	PartOfSpeech string   `json:"partOfSpeech"`
	Explain      string   `json:"explain"`
	Sentences    []string `json:"sentences"`
}

var once sync.Once
var template string

func (g Group) Ask(ctx *gin.Context) {
	var request struct {
		UserId string `json:"userId"`
		Word   string `json:"word"`
	}

	var response struct {
		api.BaseResponse
		Answer Answer `json:"answer"`
	}
	response.Answer = Answer{
		Vocs: []Voc{},
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ArgumentError, err.Error()))
		return
	}

	response.Answer.Word = request.Word

	once.Do(func() {
		example := []Voc{
			{
				Word:         "bark",
				PartOfSpeech: "vt",
				Explain:      "to shout or speak loudly and insistently",
				Sentences: []string{
					"The dog barked at the intruder.", "The coach barked orders at the players.", "He barked out a command to stop.",
				},
			},
		}

		t, _ := json.Marshal(example)
		template = string(t)
	})

	content := fmt.Sprintf(`對以下英文單字做解釋並提供3個例句:%s。以json方式輸出，key包含partOfSpeech,explain,sentences。範例: %s`, request.Word, template)

	reply, err := pixie.Chat([]openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}

	if err := json.Unmarshal([]byte(reply), &response.Answer.Vocs); err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response)
	return
}
