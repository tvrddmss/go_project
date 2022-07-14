package rpc

import (
	"encoding/gob"
	"fmt"
	"go_project/models"
)

// 定义用户对象
type UserLoginInfo struct {
	Name string
	Ip   string
}

// 用于写入日至的方法
func appendLog(uid UserLoginInfo) error {
	fmt.Println("log-username", uid.Name)
	models.AddAuth_log(uid.Name, uid.Ip)
	return nil
}

func InitLogServer() {
	// 编码中有一个字段是interface{}时，要注册一下
	gob.Register(UserLoginInfo{})
	addr := "127.0.0.1:8001"
	// 创建服务端
	srv := NewServer(addr)
	// 将服务端方法，注册一下
	srv.Register("appendLog", appendLog)
	// 服务端等待调用
	go srv.Run()
}
