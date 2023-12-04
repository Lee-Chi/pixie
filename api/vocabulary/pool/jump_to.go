package pool

import (
	"net/http"
	"pixie/api"
	"pixie/core/vocabulary"
	"pixie/db"
	"pixie/db/model"

	"github.com/gin-gonic/gin"
)

func (g Group) JumpTo(ctx *gin.Context) {
	type Definition struct {
		Text         string `json:"text"`
		PartOfSpeech string `json:"partOfSpeech"`
	}
	var request struct {
		UserId string `json:"userId"`
		Word   string `json:"word"`
	}

	var response struct {
		api.BaseResponse
		Ok          bool         `json:"ok"`
		Id          string       `json:"id"`
		Word        string       `json:"word"`
		HasToggled  bool         `json:"hasToggled"`
		Definitions []Definition `json:"definitions"`
	}
	response.Definitions = []Definition{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ArgumentError, err.Error()))
		return
	}

	voc, err := vocabulary.JumpTo(ctx, request.UserId, request.Word)
	if err != nil {
		response.Ok = false
		ctx.JSON(http.StatusOK, response)
		return
	}

	count, err := db.Pixie().Collection(model.CVocabularyBookmark).Count(
		ctx,
		model.Field_VocabularyBookmark.UserId.Equal(request.UserId).And(model.Field_VocabularyBookmark.VocabularyId.Equal(voc.Id)),
	)
	if err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}

	response.Ok = true
	response.Id = voc.Id.Hex()
	response.Word = voc.Word
	response.HasToggled = count > 0
	for _, def := range voc.Definitions {
		response.Definitions = append(response.Definitions, Definition{
			Text:         def.Text,
			PartOfSpeech: def.PartOfSpeech,
		})
	}

	ctx.JSON(http.StatusOK, response)
	return
}
