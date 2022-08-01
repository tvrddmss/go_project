package api

import (

	//"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

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
// @Router /auth/login [get]
func AuthLogin(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	userip := c.Request.RemoteAddr
	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		userid := models.Login(username, password)
		if userid > 0 {
			userinfo := models.GetUserInfo(userid)
			token, err := redis.Get("userid:" + strconv.Itoa(userid) + ":token")
			fmt.Println("token:", token)
			istoken := false
			if err == nil && token != nil {
				istoken = true
			} else {
				fmt.Println("newtoken:start")
				token, err = util.GenerateToken(username, c.Request.RemoteAddr)
				if err != nil {
					code = e.ERROR_AUTH_TOKEN
				} else {
					istoken = true
					redis.Set("user:"+strconv.Itoa(userid)+":token", token)
				}
			}
			if istoken {
				data["token"] = token
				fmt.Println("token", token)
				data["id"] = userinfo.ID
				fmt.Println("id", userinfo.ID)
				data["username"] = userinfo.Username
				fmt.Println("username", userinfo.Username)
				data["nickname"] = userinfo.Nickname
				fmt.Println("nickname", userinfo.Nickname)
				data["img"] = userinfo.Img
				fmt.Println("img", userinfo.Img)
				code = e.SUCCESS
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

// @Summary 获取用户信息
// @Tags Auth
// @Produce  json
// @param token header string true "Authorization"
// @Param id path int true "ID"
// @Success 200 {string} e.BackStruct "{"code":200,"data":{},"msg":"ok"}"
// @Router /auth/{id} [get]
func AuthGetAuth(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetUserInfo(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// @Summary 修改用户信息
// @Tags Auth
// @Produce  json
// @param token header string true "Authorization"
// @Param id path int true "ID"
// @Param nickname formData string false "nickname"
// @Success 200 {string} e.BackStruct "{"code":200,"data":{},"msg":"ok"}"
// @Router /auth/{id} [put]
func AuthEditAuth(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	//nickname := c.Query("nickname")
	nickname, _ := c.GetPostForm("nickname")

	valid := validation.Validation{}

	valid.Required(id, "id").Message("ID不能为空")
	valid.MaxSize(nickname, 100, "name").Message("昵称最长为100字符")

	var resultdata interface{}
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistAuthByID(id) {
			data := make(map[string]interface{})
			data["nickname"] = nickname
			models.EditAuth(id, data)

			resultdata = models.GetUserInfo(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": resultdata,
	})
}

// @Summary 清空缓存
// @Tags Auth
// @Produce  json
// @Success 200 {string} e.BackStruct "{"code":200,"data":{},"msg":"ok"}"
// @Router /clearredis [get]
func AuthClearRedis(c *gin.Context) {

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

//
// @Summary 上传图片
// @Tags Auth
// @Produce  json
// @param token header string true "Authorization"
// @Param userid query string false "userid"
// @Success 200 {string} e.BackStruct "{"code":200,"data":{},"msg":"ok"}"
// @Router /auth/uploadimg [post]
func AuthUploadImg(c *gin.Context) {

	baseurl := "http://192.168.50.8:3001/file/userimgs/"
	basedir := "/home/wsl/project/gowork/src/go_porject_nuxt/assets/file/userimgs/"
	//config:=config.CreateConfig()
	userid, _ := c.GetPostForm("id")
	f, err := c.FormFile("imgfile")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	} else {

		fileExt := strings.ToLower(path.Ext(f.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!只允许png,jpg,gif,jpeg文件",
			})
			return
		}
		//fileName := tools.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fileName := userid + "-" + time.Now().Format("2006-01-02 15:04:05")
		fildDir := basedir
		// isExist, _ := tools.IsFileExist(fildDir)
		// if !isExist {
		// 	os.Mkdir(fildDir, os.ModePerm)
		// }
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)
		int_userid, _ := strconv.Atoi(userid)
		data := models.GetUserInfo(int_userid)
		if len(data.Img) > 0 {
			oldpath := strings.Replace(data.Img, baseurl, basedir, -1)
			os.Remove(oldpath)
		}
		data.Img = strings.Replace(filepath, basedir, baseurl, -1)
		models.EditAuth(int_userid, data)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": filepath,
			},
		})
	}
}
