package main

import ( 
	"context"
	"log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
	//"github.com/gin-gonic/gin"
    "go_project/app/shop"
    "fmt"
    "go_project/routers"
    "go_project/pkg/setting"
)
// @title go_project学习项目
// @version 1.0
// @description 测试用程序
// @termsOfService [http://swagger.io/terms/](http://swagger.io/terms/)
// @contact.name 这里写联系人信息
// @contact.url [http://www.swagger.io/support](http://www.swagger.io/support)
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url [http://www.apache.org/licenses/LICENSE-2.0.html](http://www.apache.org/licenses/LICENSE-2.0.html)
// @host 192.168.50.237
// @BasePath /
func main() {
    // 加载多个APP的路由配置
    routers.Include(shop.Routers)
    // 初始化路由
    r := routers.Init()
    // if err := r.Run(); err != nil {
    //     fmt.Println("startup service failed, err:%v\n", err)
    // }
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