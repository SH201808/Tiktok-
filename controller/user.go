package controller

import (
	"fmt"
	"net/http"
	"tiktok/dao/mysql"
	"tiktok/models"
	"tiktok/pkg"
	"tiktok/utils/tokenUtils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//输入是否规范
	if ok := pkg.Check(username, password); !ok {
		models.ResponseErr(c, fmt.Sprintf("输入不规范"))
		return
	}
	//用户名是否冲突
	err := mysql.DB.Where("username = ?", username).Error
	if err != nil {
		models.ResponseErr(c, fmt.Sprintf("用户名已存在"))
		return
	}
	//创建用户
	user := &models.User{
		Username: username,
		Password: password,
	}
	mysql.DB.Create(user)
	//获得用户ID
	err = mysql.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		models.ResponseErr(c, fmt.Sprintf("数据库连接错误"))
		return
	}
	//生成token返回
	token, err := tokenUtils.CreateToken(user.ID)
	c.JSON(http.StatusOK, &models.TokenResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		Token:  token,
		UserId: int64(user.ID),
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//输入是否规范
	if ok := pkg.Check(username, password); !ok {
		models.ResponseErr(c, fmt.Sprintf("输入不规范"))
		return
	}
	//判断用户名和密码是否错误
	user := &models.User{}
	err := mysql.DB.Where("username = ? and password = ?", username, password).First(&user).Error
	if err != nil {
		models.ResponseErr(c, fmt.Sprintf("用户名或密码错误"))
		return
	}
	//生成token并返回
	token, err := tokenUtils.CreateToken(user.ID)
	c.JSON(http.StatusOK, &models.TokenResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		Token:  token,
		UserId: int64(user.ID),
	})
}
