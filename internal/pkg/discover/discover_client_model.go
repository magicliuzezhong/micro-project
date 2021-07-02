//
// Package discover
// @Description：服务发现model
// @Author：liuzezhong 2021/6/28 10:55 上午
// @Company cloud-ark.com
//
package discover

import (
	"github.com/hashicorp/consul/api"
	"sync"
)

//
// ConsulDiscoveryClientInstance
// @Description: 服务发现客户端实例
//
type ConsulDiscoveryClientInstance struct {
	// 客户端主机
	Host string
	//客户端端口
	Port int
	// 连接 consul 的配置
	config *api.Config
	// consul客户端
	client *api.Client
	// 互斥锁
	mutex sync.Mutex
	// 服务实例缓存字段
	instancesMap sync.Map
}
