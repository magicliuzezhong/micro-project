//
// Package main
// @Description：
// @Author：liuzezhong 2021/7/1 7:16 下午
// @Company cloud-ark.com
//
package main

import (
	"fmt"
	"time"
)

var exitExploreChannel = make(chan bool, 1)
var exitExploreChannel1 = make(chan bool, 1)

func mysqlExplore() {
	for {
		select {
		case <-exitExploreChannel:
			<-exitExploreChannel1
		case <-time.After(time.Second * 3):
			fmt.Println("程序执行")
		}
	}
}

func testFun() {
	fmt.Println("程序执行1")
	exitExploreChannel <- true
	fmt.Println("程序执行2")
	time.Sleep(time.Second * 10)
	exitExploreChannel1 <- true
	fmt.Println("程序执行3")
}

func main() {
	go mysqlExplore()
	fmt.Println("执行休眠中")
	time.Sleep(time.Second * 10)
	testFun()
	time.Sleep(time.Second * 10)
	testFun()
}
