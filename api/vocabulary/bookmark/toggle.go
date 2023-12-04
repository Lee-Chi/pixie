package bookmark

import (
	"net/http"
	"pixie/api"
	"pixie/core/vocabulary"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (g Group) Toggle(ctx *gin.Context) {
	var request struct {
		UserId       string `json:"userId"`
		VocabularyId string `json:"vocabularyId"`
	}

	var response struct {
		api.BaseResponse
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

	vocabularyId, err := primitive.ObjectIDFromHex(request.VocabularyId)
	if err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}

	if err := vocabulary.Toggle(ctx, userId, vocabularyId); err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response)
	return
}
