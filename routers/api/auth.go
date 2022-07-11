package api

import (
    "time"
	//"log"
	"net/http"

    "fmt"
	"github.com/gin-gonic/gin"
	"github.com/astaxie/beego/validation"

	"go_project/pkg/e"
	"go_project/pkg/util"
	"go_project/models"
	"go_project/pkg/logging"

	"github.com/go-redis/redis"
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

	redisdb := redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379", // redis地址
			Password: "", // redis没密码，没有设置，则留空
			DB:       0,  // 使用默认数据库
		})
	defer redisdb.Close()

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
			token,err := redisdb.Get("user:" + username + ":token").Result()
			fmt.Println("token:%s",token)
			if err == nil {
				data["token"] = token
				code = e.SUCCESS
			} else {
				
				token, err := util.GenerateToken(username, password)
				if err != nil {
					code = e.ERROR_AUTH_TOKEN
				} else {
					fmt.Println("newtoken:%s",token)
					data["token"] = token

					code = e.SUCCESS

					redisdb.Set("user:" + username + ":token",token,time.Minute*10)
				}
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

// @Summary 清空缓存
// @Tags Auth
// @Produce  json
// @Success 200 {string} e.BackStruct "{"code":200,"data":{},"msg":"ok"}"
// @Router /clearredis [get]
func ClearRedis(c *gin.Context) {

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})

	redisdb := redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379", // redis地址
			Password: "", // redis没密码，没有设置，则留空
			DB:       0,  // 使用默认数据库
		})
	defer redisdb.Close()

	var cursor uint64
	keys,cursor,err := redisdb.Scan(cursor,"*",100).Result()
	if err !=nil{
		fmt.Println("scan keys failed err:",err)
		code = e.INVALID_PARAMS
	}
	for _,key := range keys{
		redisdb.Do("del",key)
		//fmt.Println("key:",key)
		// sType,err := redisdb.Type(ctx,key).Result()
		// if err !=nil{
		// 	fmt.Println("get type failed :",err)
		// 	return
		// }
		// fmt.Printf("key :%v ,type is %v\n",key,sType)
		// if sType == "string" {
		// 	val,err := redisdb.Get(ctx,key).Result()
		// 	if err != nil{
		// 		fmt.Println("get key values failed err:",err)
		// 		return
		// 	}
		// 	fmt.Printf("key :%v ,value :%v\n",key,val)
		// }else if sType == "list"{
		// 	val,err := redisdb.LPop(ctx,key).Result()
		// 	if err !=nil{
		// 		fmt.Println("get list value failed :",err)
		// 		return
		// 	}
		// 	fmt.Printf("key:%v value:%v\n",key,val)
		// }
	}
	code = e.SUCCESS

	c.JSON(http.StatusOK, gin.H{
        "code" : code,
        "msg" : e.GetMsg(code),
        "data" : data,
    })
}