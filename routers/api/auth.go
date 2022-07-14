package api

import (

	//"log"
	"net/http"

	"fmt"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"go_project/models"
	"go_project/pkg/e"
	"go_project/pkg/logging"
	"go_project/pkg/util"
	"go_project/rpc"

	"go_project/middleware/redis"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary 登录
// @Tags Auth
// @Produce  json
// @Param username query string false "username"
// @Param password query string false "password"
// @Success 200 {string} e.BackStruct "{"code":200,"data":{},"msg":"ok"}"
// @Router /auth [get]
func GetAuth(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	userip := c.Request.RemoteAddr
	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {

			token, err := redis.Get("user:" + username + ":token")
			fmt.Println("token:", token)
			if err == nil && token != nil {
				data["token"] = token
				code = e.SUCCESS
			} else {
				fmt.Println("newtoken:start")
				token, err := util.GenerateToken(username, c.Request.RemoteAddr)
				if err != nil {
					code = e.ERROR_AUTH_TOKEN
				} else {
					fmt.Println("newtoken:", token)
					data["token"] = token

					code = e.SUCCESS

					redis.Set("user:"+username+":token", token)
				}
			}

			rpc.Send(username, userip)

		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			//log.Println(err.Key, err.Message)
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// @Summary 清空缓存
// @Tags Auth
// @Produce  json
// @Success 200 {string} e.BackStruct "{"code":200,"data":{},"msg":"ok"}"
// @Router /clearredis [get]
func ClearRedis(c *gin.Context) {

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})

	redis.Clear()
	err := recover()
	if err != nil {
		code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
	} else {
		code = e.SUCCESS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
