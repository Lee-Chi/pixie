package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"pixie/core/user"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ErrorCode_Ok int = iota
	ErrorCode_ApiError
	ErrorCode_ArgumentError
	ErrorCode_AuthError
)

type BaseResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func PackErrorResponse(code int, message string) *BaseResponse {
	return &BaseResponse{
		ErrorCode:    code,
		ErrorMessage: message,
	}
}

func CheckAuth(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")

	datas := strings.Split(auth, ".")
	if len(datas) != 3 {
		ctx.JSON(http.StatusOK, PackErrorResponse(ErrorCode_AuthError, "invalid header"))
		ctx.Abort()
		return
	}

	userId, err := primitive.ObjectIDFromHex(datas[0])
	if err != nil {
		ctx.JSON(http.StatusOK, PackErrorResponse(ErrorCode_AuthError, "wrong user id"))
		ctx.Abort()
		return
	}

	timestamp := datas[1]

	actual := datas[2]

	session, err := user.GetSession(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, PackErrorResponse(ErrorCode_AuthError, err.Error()))
		ctx.Abort()
		return
	}

	hash := sha256.New()
	hash.Write([]byte(strings.Join([]string{timestamp, session.Id.Hex()}, "")))
	expect := fmt.Sprintf("%x", hash.Sum(nil))
	if actual != expect {
		ctx.JSON(http.StatusOK, PackErrorResponse(ErrorCode_AuthError, fmt.Sprintf("actual: %s, expect: %s", actual, expect)))
		ctx.Abort()
		return
	}

	ctx.Next()
}
