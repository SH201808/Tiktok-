package routers

import (
	"tiktok/controller"

	"github.com/gin-gonic/gin"
)

func SetUpRoute() *gin.Engine {
	r := gin.Default()
	r.POST("/user/register", controller.Register)
	return r
}
