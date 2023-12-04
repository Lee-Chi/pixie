package member

import (
	"net/http"
	"pixie/api"
	"pixie/core/user"

	"github.com/gin-gonic/gin"
)

func (g Group) Login(ctx *gin.Context) {
	var request struct {
		Account  string `json:"account"`
		Password string `json:"password"`
		Device   string `json:"device"`
	}

	var response struct {
		api.BaseResponse
		UserId    string `json:"userId"`
		SessionId string `json:"sessionId"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ArgumentError, err.Error()))
		return
	}

	u, err := user.FindOneMember(ctx, request.Account)
	if err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}
	response.UserId = u.Id.Hex()

	sessionId, err := u.Login(ctx, request.Password, request.Device)
	if err != nil {
		ctx.JSON(http.StatusOK, api.PackErrorResponse(api.ErrorCode_ApiError, err.Error()))
		return
	}

	response.SessionId = sessionId

	ctx.JSON(http.StatusOK, response)
	return
}
