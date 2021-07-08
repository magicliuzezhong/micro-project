//
// Package discover
// @Description：服务注册与发现consul实现
// @Author：liuzezhong 2021/6/28 10:57 上午
// @Company cloud-ark.com
//
package discover

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hashicorp/go-uuid"
	"micro-project/internal/pkg/common"
	"micro-project/internal/pkg/util"
	"strconv"
)

type ConsulDiscoverClient struct {
	ConsulDiscoveryClientInstance
	InstanceId string
}

func NewConsulDiscoverClient(host string, port int) *ConsulDiscoverClient {
	// 通过 Consul Host 和 Consul Port 创建一个 consul.Client
	var consulConfig = api.DefaultConfig()
	consulConfig.Address = host + ":" + strconv.Itoa(port)
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil
	}
	uuid, _ := uuid.GenerateUUID()
	return &ConsulDiscoverClient{
		ConsulDiscoveryClientInstance: ConsulDiscoveryClientInstance{
			Host:   host,
			Port:   port,
			config: consulConfig,
			client: apiClient,
		},
		InstanceId: uuid,
	}
}

func (consulClient *ConsulDiscoverClient) Register(svcHost string, port int, svcName string,
	weight int, meta map[string]string, tags []string) bool {
	var application = util.GetApplication()
	var checkUrl = fmt.Sprintf("http://%s:%d/health", util.GetLocalIp(), application.Server.Port)
	// 1. 构建服务实例元数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      consulClient.InstanceId,
		Name:    svcName,
		Address: svcHost,
		Port:    port,
		Meta:    meta,
		Tags:    tags,
		Weights: &api.AgentWeights{
			Passing: weight,
		},
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           checkUrl,
			Interval:                       "5s",
		},
	}
	var err = consulClient.client.Agent().ServiceRegister(serviceRegistration)
	if err != nil {
		util.GetLogger().Errorln("注册到注册中心失败，", err.Error())
		return false
	}
	util.GetLogger().Println("注册到注册中心成功")
	return true
}

func (consulClient *ConsulDiscoverClient) DeRegister() bool {
	var err = consulClient.client.Agent().ServiceDeregister(consulClient.InstanceId)
	if err != nil {
		util.GetLogger().Errorln("注销注册中心失败，", err.Error())
		return false
	}
	util.GetLogger().Println("注销注册中心成功")
	return true
}

func (consulClient *ConsulDiscoverClient) DiscoverServices(serviceName string) []*common.ServiceInstance {
	//  该服务已监控并缓存
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}
	// 申请锁
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()
	// 再次检查是否监控
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	} else {
		// 注册监控
		go func() {
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return // 数据异常，忽略
				}
				// 没有服务实例在线
				if len(v) == 0 {
					util.GetLogger().Println("当前没有服务实例在线")
					consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
				}
				var healthServices []*common.ServiceInstance
				for _, service := range v {
					if service.Checks.AggregatedStatus() == api.HealthPassing {
						healthServices = append(healthServices, parseAgent(service.Service))
					}
				}
				consulClient.instancesMap.Store(serviceName, healthServices)
			}
			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	}
	// 根据服务名请求服务实例列表
	entries, _, err := consulClient.client.Catalog().Service(serviceName, "", nil)
	if err != nil {
		consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
		util.GetLogger().Println("发现服务实例异常，", err.Error())
		return nil
	}
	instances := make([]*common.ServiceInstance, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = parseCatalog(entries[i])
	}
	consulClient.instancesMap.Store(serviceName, instances)
	return instances
}

//
// parseCatalog
// @Description: 解析catalog
// @param service catalog
// @return *common.ServiceInstance 服务实例
//
func parseCatalog(service *api.CatalogService) *common.ServiceInstance {
	return newServiceInstance(service.Address, service.ServicePort, service.ServiceMeta, service.ServiceWeights.Passing)
}

//
// parseAgent
// @Description: 解析agent
// @param agent agent
// @return *common.ServiceInstance 服务实例
//
func parseAgent(agent *api.AgentService) *common.ServiceInstance {
	return newServiceInstance(agent.Address, agent.Port, agent.Meta, agent.Weights.Passing)
}

//
// newServiceInstance
// @Description: 构建服务实例
// @param address 地址
// @param port 端口
// @param meta 元数据
// @param weight 权重
// @return *common.ServiceInstance 实例
//
func newServiceInstance(address string, port int, meta map[string]string, weight int) *common.ServiceInstance {
	var rpcPort = port - 1
	if meta != nil {
		if rpcPortString, ok := meta["rpcPort"]; ok {
			rpcPort, _ = strconv.Atoi(rpcPortString)
		}
	}
	return &common.ServiceInstance{
		Host:     address,
		Port:     port,
		GrpcPort: rpcPort,
		Weight:   weight,
	}
}
