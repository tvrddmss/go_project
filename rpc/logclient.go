package rpc

import (
	"fmt"
	"net"
)

func Send(username string, userip string) {
	//写日志
	addr := "127.0.0.1:8001"
	// 客户端获取连接
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("err")
	}
	// 创建客户端对象
	cli := NewClient(conn)
	// 需要声明函数原型
	var appendLog func(userinfo UserLoginInfo) error
	cli.callRPC("appendLog", &appendLog)
	// 得到查询结果
	var userinfo UserLoginInfo
	userinfo.Name = username
	userinfo.Ip = userip
	appendLog(userinfo)
	// if err != nil {
	// 	fmt.Println("err")
	// }
	fmt.Println("wancheng")
}
