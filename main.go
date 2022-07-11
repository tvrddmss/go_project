package main

import ( 
	"context"
	"log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
	"github.com/gin-gonic/gin"
    "go_project/app/shop"
    "fmt"
    "go_project/routers"
    "go_project/pkg/setting"


    _ "go_project/docs"
	// ginSwagger "github.com/swaggo/gin-swagger"
    // swaggerFiles "github.com/swaggo/files"

	"github.com/go-redis/redis"
)

var swagHandler gin.HandlerFunc

// func Init() {
//     swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
// }


// redis

// 定义一个全局变量
var redisdb *redis.Client

func initRedis()(err error){
	redisdb = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",  // 指定
		Password: "",
		DB:0,		// redis一共16个库，指定其中一个库即可
	})
    _,err = redisdb.Ping().Result()
	return
}


// @title go_project学习项目
// @version 1.0
// @description 测试用程序
// @termsOfService [http://swagger.io/terms/](http://swagger.io/terms/)
// @contact.name 这里写联系人信息
// @contact.url [http://www.swagger.io/support](http://www.swagger.io/support)
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url [http://www.apache.org/licenses/LICENSE-2.0.html](http://www.apache.org/licenses/LICENSE-2.0.html)
// @host 192.168.50.237:8081
// @BasePath /
func main() {

    err := initRedis()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n",err)
		return
	}
	fmt.Println("redis连接成功！")


    // 存普通string类型，10分钟过期
	redisdb.Set("test:name","科科儿子",time.Minute*10)
	// 存hash数据
	redisdb.HSet("test:class","521",42)
	// 存list数据
	redisdb.RPush("test:list",1)  // 向右边添加元素
	redisdb.LPush("test:list",2)  // 向左边添加元素
	// 存set数据
	redisdb.SAdd("test:set","apple")
	redisdb.SAdd("test:set","pear")



    // 加载多个APP的路由配置
    routers.Include(shop.Routers)
    // 初始化路由
    r := routers.Init()
    // if err := r.Run(); err != nil {
    //     fmt.Println("startup service failed, err:%v\n", err)
    // }

    if swagHandler != nil { 
       r.GET("/swagger/*any", swagHandler)        
    }
	svr := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	//svr.ListenAndServe()


	// r := gin.Default()
	//r.LoadHTMLGlob("templates/**/*")
	// r.GET("/kitty", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello Kitty",
	// 	})
	// })

	// r.GET("/posts/index", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "posts/index.html", gin.H{
	// 		"title": "posts/index",
	// 	})
	// })
    //routers.LoadShop(r)
    // routers.LoadShop(r)
    // if err := r.Run(); err != nil {
    //     fmt.Println("startup service failed, err:%v\n", err)
    // }

	// r.GET("/notify/signal/shutdown", func(c *gin.Context) {
    //     time.Sleep(5 * time.Second)
    //     c.String(http.StatusOK, "Test notify signal to shutdown server !")
    // })

    // svr := http.Server{
    //     Addr:    ":8080",
    //     Handler: r,
    // }

    go func() {
        // 开启一个goroutine启动服务
        if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
    quit := make(chan os.Signal, 1) // 创建一个接收信号的通道

    // kill 默认会发送 syscall.SIGTERM 信号
    // kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
    // kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
    // signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
    <-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行

    log.Println("Shutdown Server ...")
    // 创建一个10秒超时的context
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
    if err := svr.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown: ", err)
    }

    log.Println("Server exited")


	//r.Run()
}