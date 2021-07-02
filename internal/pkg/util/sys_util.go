//
// Package util
// @Description：系统工具类
// @Author：liuzezhong 2021/6/25 5:46 下午
// @Company cloud-ark.com
//
package util

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

//
// ListenExitSignal
// @Description: 监听退出信号
// @param callback 回调函数
//
func ListenExitSignal(callback func()) {
	var errChannel = make(chan error)
	go func() {
		GetLogger().Println("监听linux的ctrl c命令和kill -9命令")
		var exitChannel = make(chan os.Signal)
		signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)
		errChannel <- fmt.Errorf("%d", <-exitChannel)
	}()
	GetLogger().Errorln("进程退出，退出信号：", (<-errChannel).Error())
	if callback != nil { //回调函数非空，执行回调函数，执行程序退出方法
		callback()
	}
}

//
// GetLocalIp
// @Description: 获取本地IP
// @return string IP
//
func GetLocalIp() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addresses { // 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
