package sentence

import (
	"fmt"
	"net/http"
	"pixie/api"
	"pixie/core/pixie"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func Explain(ctx *gin.Context) {
	var request struct {
		Text string `json:"text"`
	}

	var response struct {
		api.BaseResponse
		Explained string `json:"explained"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ArgumentError, err.Error()))
		return
	}

	content := fmt.Sprintf(`對以下文字解釋和分析文法: %s`, request.Text)

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

	response.Explained = reply

	ctx.JSON(http.StatusOK, response)
	return
}
