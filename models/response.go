package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type TokenResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

func ResponseErr(c *gin.Context, msg string) {
	rE := &Response{
		StatusCode: 1,
		StatusMsg:  msg,
	}
	c.JSON(http.StatusOK, rE)
}
