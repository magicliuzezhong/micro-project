//
// Package main
// @Description：
// @Author：liuzezhong 2021/6/28 4:38 下午
// @Company cloud-ark.com
//
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	services "micro-project/internal/app/rpc/service"
)

func main() {
	var conn, err = grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		fmt.Println("服务异常", err.Error())
		return
	}
	defer conn.Close()
	var productClient = services.NewProdServiceClient(conn)
	proResponse, err := productClient.ProService(context.Background(), &services.ProdRequest{
		ProId: 26,
	})
	if err != nil {
		fmt.Println("请求异常")
		return
	}
	fmt.Println(proResponse)
}
