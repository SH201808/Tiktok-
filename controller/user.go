package controller

import (
	"net/http"
	"tiktok/pkg"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if ok := pkg.Check(username, password); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
		})
		return
	}

}
