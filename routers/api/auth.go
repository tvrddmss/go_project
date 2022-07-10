package api

import (
	//"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/astaxie/beego/validation"

	"go_project/pkg/e"
	"go_project/pkg/util"
	"go_project/models"
	"go_project/pkg/logging"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}


// @Summary 登录
// @Produce  json
// @Param username query string false "username"
// @Param password query string false "password"
// @Success 200 {string} e.BackStruct "{"code":200,"data":{},"msg":"ok"}"
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token

				code = e.SUCCESS
			}

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
        "code" : code,
        "msg" : e.GetMsg(code),
        "data" : data,
    })
}