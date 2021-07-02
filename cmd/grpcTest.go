//
// Package main
// @Description：
// @Author：liuzezhong 2021/6/28 4:33 下午
// @Company cloud-ark.com
//
package main

import (
	"google.golang.org/grpc"
	services "micro-project/internal/app/rpc/service"
	"net"
)

func main() {
	var grpcServer = grpc.NewServer()
	services.RegisterProdServiceServer(grpcServer, new(services.ProductService))
	var listen, _ = net.Listen("tcp", ":8081")
	grpcServer.Serve(listen)
}
