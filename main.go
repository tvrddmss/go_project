package main

import (
	"context"
	"fmt"
	"go_project/app/shop"
	"go_project/middleware/channel"
	_ "go_project/middleware/channel"
	"go_project/middleware/redis"
	"go_project/pkg/setting"
	"go_project/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	_ "go_project/docs"
	"go_project/rpc"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// swaggerFiles "github.com/swaggo/files"
)

var swagHandler gin.HandlerFunc

// func Init() {
//     swagHandler = ginSwagger.WrapHandler(swaggerFiles.Handler)
// }

// @title go_project学习项目
// @version 1.0
// @description 测试用程序
// @termsOfService [http://swagger.io/terms/](http://swagger.io/terms/)
// @contact.name 这里写联系人信息
// @contact.url [http://www.swagger.io/support](http://www.swagger.io/support)
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url [http://www.apache.org/licenses/LICENSE-2.0.html](http://www.apache.org/licenses/LICENSE-2.0.html)
// @host 192.168.50.8:8000
// @BasePath /
func main() {

	//initChannel()
	//启动微服务
	rpc.InitLogServer()
	initaaa()

}
func initaaa() {
	redis.Init()
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

var channel1 chan int

func initChannel() {
	channel1 = make(chan int)
	fmt.Println("初始化完毕")
	go ChannelRescive()
	go ChannelSend()
	time.Sleep(time.Second)
}

func ChannelSend() {
	channel.Send(10, channel1)
	fmt.Println("发送：10")
}
func ChannelRescive() {
	fmt.Println("接收开始")
	a := channel.Resvice(channel1)
	fmt.Printf("接收：%d\n", a)
}
