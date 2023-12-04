package sentence

import (
	"fmt"
	"net/http"
	"pixie/api"
	"pixie/core/pixie"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func Interpret(ctx *gin.Context) {
	var request struct {
		Sentence string `json:"sentence"`
	}

	var response struct {
		api.BaseResponse
		Interpreted string `json:"interpreted"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ArgumentError, err.Error()))
		return
	}

	content := fmt.Sprintf(`翻譯以下英文句子為繁體中文: %s`, request.Sentence)

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

	response.Interpreted = reply

	ctx.JSON(http.StatusOK, response)
	return
}
