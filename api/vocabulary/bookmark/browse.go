package bookmark

import (
	"net/http"
	"pixie/api"
	"pixie/core/vocabulary"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (g Group) Browse(ctx *gin.Context) {
	var request struct {
		UserId string `json:"userId"`
		Page   int64  `json:"page"`
	}

	type Voc struct {
		Id           string `json:"id"`
		Word         string `json:"word"`
		PartOfSpeech string `json:"partOfSpeech"`
		Text         string `json:"text"`
	}
	var response struct {
		api.BaseResponse
		Vocs []Voc `json:"vocs"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ArgumentError, err.Error()))
		return
	}

	userId, err := primitive.ObjectIDFromHex(request.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}

	page := request.Page
	if page < 1 {
		page = 1
	}
	vocs, err := vocabulary.BrowseBookmark(ctx, userId, page)
	if err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}

	response.Vocs = []Voc{}
	for _, voc := range vocs {
		for _, def := range voc.Definitions {
			response.Vocs = append(response.Vocs, Voc{
				Id:           voc.Id.Hex(),
				Word:         voc.Word,
				PartOfSpeech: def.PartOfSpeech,
				Text:         def.Text,
			})
		}
	}

	ctx.JSON(http.StatusOK, response)
	return
}
